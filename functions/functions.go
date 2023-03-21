package functions

import (
	"fmt"
	"math/rand"
)

func Hello() {
	fmt.Println("Hello, World!")
}

type creature struct {
	Hole    int16
	Health  int16
	Respect int16
	Weight  int16
}

func NewPlayer() creature {
	return creature{
		Hole:    8,
		Health:  100,
		Respect: 10,
		Weight:  20,
	}
}

func (c *creature) CheckWin() bool {
	if c.Respect >= 100 {
		return true
	} else {
		return false
	}
}

func (c *creature) CheckDefeat() bool {
	if c.Health <= 0 || c.Hole <= 0 || c.Respect <= 0 || c.Weight <= 0 {
		return true
	} else {
		return false
	}
}

func (c *creature) Info() {
	fmt.Printf("Player's characteristics: Hole %v Health  %v Respect  %v Weight  %v\n",
		c.Hole, c.Health, c.Respect, c.Weight)
}

func (c *creature) Day() {
	var input int8
	fmt.Print("\tThe day has come..\n")
	fmt.Print("Choose the action:\n" +
		"1. Dig the hole\n" +
		"2. Eat a grass\n" +
		"3. Fight\n" +
		"4. Sleep\n" +
		"Your choice: ")
	fmt.Scan(&input)
	switch input {
	case 1:
		c.dig()
	case 2:
		c.eat()
	case 3:
		c.fight()
	case 4:
		c.sleep()
	}
}

func (c *creature) Night() {
	fmt.Print("You fall asleep\n")
	c.Hole -= 2
	c.Health += 20
	c.Respect -= 2
	c.Weight -= 5
	// for i := 1; i <= 3; i++ {
	// 	fmt.Print(".")
	// 	time.Sleep(1 * time.Second)
	// }
	c.Info()
}

func (c *creature) dig() {
	var input int8
	fmt.Print("How will you dig:\n" +
		"1. High\n" +
		"2. Lazily\n" +
		"Your choice: ")
	fmt.Scan(&input)
	switch input {
	case 1:
		c.Hole += 5
		c.Health -= 30
	case 2:
		c.Hole += 2
		c.Health -= 10
	}
	c.Info()
}

func (c *creature) fight() {
	var input int8
	var winrate, enemyWeight, exodus float32
	fmt.Print("Who do you want to fight:\n" +
		"1. Weak (Weight 30)\n" +
		"2. Medium (Weight 50)\n" +
		"3. Strong (Weight 70)\n" +
		"Your choice: ")
	fmt.Scan(&input)
	switch input {
	case 1:
		enemyWeight = 30
	case 2:
		enemyWeight = 50
	case 3:
		enemyWeight = 70
	}
	winrate = float32(c.Weight) / (float32(c.Weight) + enemyWeight - 25)
	fmt.Printf("Winrate is %.2f percent\n", winrate*100)
	if winrate*100 < 50 {
		fmt.Print("You really want to fight?\n" +
			"1. Yes\n" +
			"2. No, i would do something else\n" +
			"Your choice: ")
		fmt.Scan(&input)
		switch input {
		case 1:
			break
		case 2:
			c.Day()
			return
		}
	}
	fmt.Println("\tThe fight is beginning!")
	exodus = rand.Float32()
	//fmt.Printf("exodus %v", exodus)
	difference := int16(enemyWeight) - c.Weight
	if winrate > exodus {
		fmt.Println("You won!")
		if c.Weight < int16(enemyWeight) {
			c.Respect += difference + 10
		}
		if c.Weight == int16(enemyWeight) {
			c.Respect += 15
		}
		if c.Weight > int16(enemyWeight) {
			c.Respect += 10
		}
	} else {
		fmt.Println("You lose")
		if c.Weight < int16(enemyWeight) {
			c.Health -= difference
		} else {
			fmt.Print("You're too strong to take damage")
		}
	}
	c.Info()
}

func (c *creature) eat() {
	var input int8
	fmt.Print("What grass will you eat:\n" +
		"1. Withered\n" +
		"2. Green\n" +
		"Your choice: ")
	fmt.Scan(&input)
	switch input {
	case 1:
		c.Weight += 15
		c.Health += 10
	case 2:
		if c.Respect >= 30 {
			c.Health += 30
			c.Weight += 30
		} else {
			c.Health -= 30
		}
	}
	c.Info()
}

func (c *creature) sleep() {
	fmt.Print("You decide to take a nap")
	c.Night()
}
