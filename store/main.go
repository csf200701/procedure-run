package store

import (
	"encoding/xml"
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"procedure-run/connector"
	"sync"
	"time"
)

var dataFile string = "pr-data.xml"
var pr *procedureRun = new(procedureRun)
var isUpdate bool = false

type procedureRun struct {
	XMLName   xml.Name  `xml:"ProcedureRun"`
	Databases DataBases `xml:"Databases"`
	lock      sync.RWMutex
}

type DataBases struct {
	Databases []DataBase `xml:"Database"`
}

type DataBase struct {
	XMLName            xml.Name              `xml:"Database"`
	DatabaseAlias      string                `xml:"DatabaseAlias"`
	DatabaseCollection *connector.Collection `xml:"DatabaseCollection"`
}

func init() {
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		if _, err := os.Create(dataFile); err != nil {
			panic(err)
		}
	}
	xmlFile, err := os.Open(dataFile)
	if err != nil {
		panic(err)
	}

	defer xmlFile.Close()
	//从xml文档中读取数据
	xmlData, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		panic(err)
	}
	xml.Unmarshal(xmlData, pr)

	startStore()
}

func FindAll() []DataBase {
	return pr.Databases.Databases
}

func AddDatabase(databaseAlias string, databaseCollection *connector.Collection) {
	pr.lock.Lock()
	defer pr.lock.Unlock()
	database, err := GetDatabase(databaseAlias)
	if err != nil {
		database = DataBase{DatabaseAlias: databaseAlias, DatabaseCollection: databaseCollection}
	}
	pr.Databases.Databases = append(pr.Databases.Databases, database)
	isUpdate = true
}

func DeleteDatabase(databaseAlias string) {
	pr.lock.Lock()
	defer pr.lock.Unlock()
	for i, database := range pr.Databases.Databases {
		if database.DatabaseAlias == databaseAlias {
			if i == 0 {
				pr.Databases.Databases = pr.Databases.Databases[1:]
			} else if i == len(pr.Databases.Databases)-1 {
				pr.Databases.Databases = pr.Databases.Databases[:i]
			} else {
				pr.Databases.Databases = append(pr.Databases.Databases[:i], pr.Databases.Databases[i+1:]...)
			}
			isUpdate = true
		}
	}
}

func GetDatabase(databaseAlias string) (DataBase, error) {
	for _, database := range pr.Databases.Databases {
		if database.DatabaseAlias == databaseAlias {
			return database, nil
		}
	}
	return DataBase{}, errors.New("不存在")
}

func startStore() {
	ticker := time.NewTicker(4 * time.Second)
	//defer ticker.Stop()
	go func() {
		for {
			<-ticker.C
			if isUpdate {
				saveXML()
				isUpdate = false
			}
		}
	}()
}

func saveXML() {
	pr.lock.Lock()
	data, err := xml.MarshalIndent(pr, " ", " ")
	if err == nil {
		ioutil.WriteFile(dataFile, data, fs.ModeAppend)
	}
	pr.lock.Unlock()
}
