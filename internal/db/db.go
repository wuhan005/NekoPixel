package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/wuhan005/NekoPixel/internal/conf"
	"github.com/wuhan005/NekoPixel/internal/dbutil"
	"github.com/wuhan005/NekoPixel/internal/strutil"
)

var AllTables = []interface{}{
	&Color{},
	&CanvasPixel{},
	&Pixel{},
}

// Init initializes the database.
func Init() (*gorm.DB, error) {
	dsn := os.ExpandEnv("postgres://$PGUSER:$PGPASSWORD@$PGHOST/$PGNAME?sslmode=$PGSSLMODE")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		NowFunc: func() time.Time {
			return dbutil.Now()
		},
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             3 * time.Second,
				LogLevel:                  logger.Silent,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		),
	})
	if err != nil {
		return nil, errors.Wrap(err, "open connection")
	}

	// Migrate databases.
	if err := db.AutoMigrate(AllTables...); err != nil {
		return nil, errors.Wrap(err, "auto migrate")
	}

	// Create sessions table.
	q := `
CREATE TABLE IF NOT EXISTS sessions (
    key        TEXT PRIMARY KEY,
    data       BYTEA NOT NULL,
    expired_at TIMESTAMP WITH TIME ZONE NOT NULL
);`
	if err := db.Exec(q).Error; err != nil {
		return nil, errors.Wrap(err, "create sessions table")
	}

	SetDatabaseStore(db)

	// Create trigger function.
	if err := db.Exec(`
CREATE OR REPLACE FUNCTION upsert_canvas_pixel()
RETURNS TRIGGER AS $$
DECLARE
    colorIndex TEXT;
BEGIN
    SELECT index INTO colorIndex FROM colors WHERE color = NEW.color LIMIT 1;

    IF colorIndex IS NOT NULL THEN
        UPDATE canvas_pixels
        SET color = NEW.color, index = colorIndex
        WHERE x = NEW.x AND y = NEW.y;

        IF NOT FOUND THEN
            INSERT INTO canvas_pixels(x, y, color, index)
            VALUES (NEW.x, NEW.y, NEW.color, colorIndex);
        END IF;
    ELSE
        RAISE EXCEPTION 'Color not found in colors table.';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
`).Error; err != nil {
		return nil, errors.Wrap(err, "create trigger function")
	}

	// Create trigger.
	if err := db.Exec(
		`CREATE OR REPLACE TRIGGER trigger_upsert_canvas_pixel
AFTER INSERT ON pixels
FOR EACH ROW
EXECUTE FUNCTION upsert_canvas_pixel();`,
	).Error; err != nil {
		return nil, errors.Wrap(err, "create trigger")
	}

	// Init colors.
	ctx := context.Background()
	isEmpty, err := Colors.IsEmpty(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "check colors")
	}
	if isEmpty {
		colors := make(map[string]string, len(conf.App.Colors))
		for idx, color := range conf.App.Colors {
			hash := strutil.GenerateCode(idx)
			colors[color] = hash
		}
		if err := Colors.Create(ctx, CreateColorOptions{Colors: colors}); err != nil {
			return nil, errors.Wrap(err, "create colors")
		}
	}

	return db, nil
}

// SetDatabaseStore sets the database table store.
func SetDatabaseStore(db *gorm.DB) {
	Colors = NewColorsStore(db)
	Pixels = NewPixelStore(db)
	Canvas = NewCanvasStore(db)
}
