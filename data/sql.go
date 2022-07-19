package data

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/jmoiron/sqlx"

	_ "modernc.org/sqlite"
)

//SQLDB 数据库初始化
var SQLDB *sqlx.DB

func Init() {
	var err error
	if !PathExist("./data") {
		os.Mkdir("./data", os.ModePerm)
	}
	SQLDB, err = sqlx.Connect("sqlite", "data/demo.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	if !PathExist("./data/install.lock") {
		err = installsql()
		if err != nil {
			log.Fatal(err.Error())
		}
		ioutil.WriteFile("./data/install.lock", []byte("install"), 0644)
	}
}
func installsql() error {
	_, err := SQLDB.Exec(`
	DROP TABLE IF EXISTS "config";
	CREATE TABLE "config" (
	  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	  "keyid" text NOT NULL DEFAULT '',
	  "keysecret" text NOT NULL DEFAULT '',
	  "ifileurl" TEXT NOT NULL DEFAULT ''
	);
	
	DROP TABLE IF EXISTS "project";
	CREATE TABLE "project" (
	  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	  "title" TEXT NOT NULL DEFAULT '',
	  "ifile_root" TEXT NOT NULL DEFAULT '',
	  "uid" integer NOT NULL DEFAULT 0
	);
	
	DROP TABLE IF EXISTS "task";
	CREATE TABLE "task" (
	  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	  "title" TEXT NOT NULL DEFAULT '',
	  "ifile_root" TEXT NOT NULL DEFAULT '',
	  "uid" INTEGER NOT NULL DEFAULT 0
	);

	DROP TABLE IF EXISTS "users";
	CREATE TABLE "users" (
	  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	  "username" TEXT NOT NULL DEFAULT '',
	  "password" TEXT NOT NULL DEFAULT '',
	  "email" TEXT NOT NULL DEFAULT '',
	  "mobile" TEXT NOT NULL DEFAULT '',
	  "ifileuid" INTEGER NOT NULL DEFAULT 0
	);
		`)
	return err
}

func PathExist(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil || os.IsExist(err)
}
