/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	//cmd.Execute()

}
