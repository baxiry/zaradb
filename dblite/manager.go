package dblite

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tidwall/gjson"
)

func NewCollection(query string) string {
	collectName := gjson.Get(query, "collection").String()

	if collectName == "" {
		return "please shose a collectin name"
	}

	// read cllections file name
	infos, err := os.OpenFile(db.Name+db.Infos, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		eLog.Println("open collectFiles : ", err)
	}
	defer infos.Close()

	names, err := io.ReadAll(infos)
	if err != nil {
		eLog.Println(err)
	}

	collectsList := strings.Split(string(names), " ")

	for _, coll := range collectsList {
		fmt.Println(coll, collectName)
		if coll == collectName {
			return collectName + " already exist!"
		}

	}

	_, err = infos.WriteString(collectName + " ")
	if err != nil {
		return fmt.Sprintf("ERROR can't create %s collectiln", collectName)
	}

	return "new collecteon is created"
}
