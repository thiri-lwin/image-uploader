package dao

import (
	"fmt"
	"go_skill_test/models"
	"go_skill_test/utilities"
)

func SaveImageInfo(imageInfo models.ImageInfo) error {

	session := utilities.DBSession.Copy()
	defer session.Close()

	//Insert data into db
	err := session.DB(utilities.DBConfig.DatabaseName).C("image_info").Insert(imageInfo)
	if err != nil {
		fmt.Println("Error in saving data into db :", err)
		return err
	}

	return nil
}
