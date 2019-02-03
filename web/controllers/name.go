package controllers

import (
	"context"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
	"github.com/yacen/song-resgain/web/models"
	"github.com/yacen/song-resgain/web/service"
	"log"
)

type NameController struct {
	beego.Controller
}

func (c *NameController) Get() {
	name:=c.GetString("name")
	page,_:=c.GetInt("page")
	fmt.Println(name,page)
	skip:=(page-1)*100
	if skip <0 {
		skip = 0
	}

	var gs []models.Girl
	ops := options.Find()
	ops.SetSort(bsonx.Doc{{"name", bsonx.Int32(1)}}).SetSkip(int64(skip)).SetLimit(100)
	cur, err := service.Collection.Find(context.TODO(),bson.D{{"explain", bson.D{{"$regex","大吉昌"}}},{"name",bson.D{{"$regex",name}}}}, ops)
	if err != nil {
		fmt.Println("========================")
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem models.Girl
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		gs = append(gs, elem)
	}
	c.Data["Items"] = gs
	if page <=1 {
		c.Data["Prev"] = 1
	}else {
		c.Data["Prev"] = page-1
	}
	c.Data["Next"] = page+1
	c.Data["Name"] = name

	c.TplName = "index.html"
}
