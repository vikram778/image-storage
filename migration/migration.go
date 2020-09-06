package migration

import (
	"image-storage/migrate"
)

// LocalMigrations ...
var LocalMigrations = migrate.Migrations{
	migrate.Migration{
		ID: 1599339905,
		SQL: `CREATE TABLE IF NOT EXISTS albums (
			  id int(11) NOT NULL AUTO_INCREMENT,
			  tittle varchar(100) DEFAULT NULL,
			  image_count int(4) NOT NULL DEFAULT '0',
			  created_at datetime DEFAULT CURRENT_TIMESTAMP,
			  updated_at datetime DEFAULT NULL,
			  PRIMARY KEY (id),
		      UNIQUE (tittle),
              KEY tittle_idx (tittle));`,
	},

	migrate.Migration{
		ID: 1599339906,
		SQL: `CREATE TABLE IF NOT EXISTS images (
			  id int(11) NOT NULL AUTO_INCREMENT,
			  album_id int(11) NOT NULL DEFAULT '0',
			  image_path varchar(400) DEFAULT NULL,
			  created_at datetime DEFAULT CURRENT_TIMESTAMP,
			  PRIMARY KEY (id),
              KEY album_id_idx (album_id));`,
	},
}
