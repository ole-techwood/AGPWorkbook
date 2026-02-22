package menu

type Pizza struct {
	ID    int
	Name  string
	Price float64
}

var Menu = []Pizza{
	{ID: 1, Name: "Margherita", Price: 8.99},
	{ID: 2, Name: "Pepperoni", Price: 9.99},
	{ID: 3, Name: "Vegetarian", Price: 8.49},
	{ID: 4, Name: "Hawaiian", Price: 10.49},
	{ID: 5, Name: "BBQ Chicken", Price: 10.99},
	{ID: 6, Name: "Quattro Formaggi", Price: 11.99},
	{ID: 7, Name: "Spicy Inferno", Price: 11.49},
}
