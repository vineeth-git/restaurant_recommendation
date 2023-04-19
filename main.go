package main

import (
	"fmt"
	"recommendations/models"
	"recommendations/strategy"
	"sort"
	"time"
)

func testt(aa string) bool {

	fmt.Println("Function called", aa)
	return false
}

func input() (models.User, []models.Restaurant) {
	user := models.User{
		Cuisines: []models.CuisineTracking{
			{"Chinese", 20},
			{"NorthIndian", 5},
			{"Biryani", 3},
		},
		CostBracket: []models.CostTracking{
			{1, 4},
			{4, 3},
			{5, 2},
		},
	}
	availableRestaurants := []models.Restaurant{
		{RestaurantID: "1", Cuisine: "Chinese", CostBracket: 5, Rating: 2, IsRecommended: true, OnboardedTime: time.Now().AddDate(0, 0, -2)},
		{RestaurantID: "2", Cuisine: "Chinese", CostBracket: 1, Rating: 4, IsRecommended: true, OnboardedTime: time.Now().AddDate(0, 0, -2)},
		{RestaurantID: "3", Cuisine: "NorthIndian", CostBracket: 3, Rating: 4, IsRecommended: true, OnboardedTime: time.Now().AddDate(0, 0, -2)},
		{RestaurantID: "4", Cuisine: "Chinese", CostBracket: 4, Rating: 4, IsRecommended: true, OnboardedTime: time.Now().AddDate(0, 0, -2)},
		{RestaurantID: "5", Cuisine: "Biryani", CostBracket: 5, Rating: 4, IsRecommended: true, OnboardedTime: time.Now().AddDate(0, 0, -2)},
	}
	return user, availableRestaurants
}

func setPrimarySecondaryOptions(user *models.User) {
	userPrimaryCuisine, userSecondaryCuisines := getPrimarySecondaryCuisine(*user)
	user.PrimaryCuisine = userPrimaryCuisine
	user.SecondaryCuisine = userSecondaryCuisines

	userPrimaryCost, userSecondaryCost := getPrimarySecondaryCost(*user)
	user.PrimaryCostBracket = userPrimaryCost
	user.SecondaryCostBracket = userSecondaryCost
}

func aggregateStrategyResults(mp map[int][]models.Restaurant) []models.Restaurant {
	uniqueRestaurant := map[models.Restaurant]bool{}
	result := []models.Restaurant{}
	resultCount := 0
	for i := 0; i < len(strategy.StrategyList); i++ {
		if i == 1 && len(result) > 0 {
			i = 3
			continue
		}
		for _, v := range mp[i] {
			if _, ok := uniqueRestaurant[v]; !ok {
				uniqueRestaurant[v] = true
				resultCount++
				result = append(result, v)
				if resultCount == 100 {
					break
				}
			}
		}
		if resultCount == 100 {
			break
		}
	}
	return result
}
func main() {
	user, availableRestaurants := input()
	setPrimarySecondaryOptions(&user)

	fmt.Println("primary cuisine - ", user.PrimaryCuisine)
	fmt.Println("secondary cuisine - ", user.SecondaryCuisine)
	fmt.Println("primary cost - ", user.PrimaryCostBracket)
	fmt.Println("secondary cost - ", user.SecondaryCostBracket)

	mp := strategy.GetRestaurantsBasedOnStrategies(user, availableRestaurants)

	result := aggregateStrategyResults(mp)
	show(result)
}

func show(ip []models.Restaurant) {
	for _, v := range ip {
		fmt.Println(v.RestaurantID, " ", v.Cuisine, " ", v.CostBracket)
	}
}

func getPrimarySecondaryCuisine(user models.User) (primary models.CuisineTracking, secondary []models.CuisineTracking) {
	if len(user.Cuisines) == 0 {
		return
	}
	sort.Slice(user.Cuisines, func(i, j int) bool {
		return user.Cuisines[i].NoOfOrders > user.Cuisines[j].NoOfOrders
	})
	primary = user.Cuisines[0]
	secondaryCuisines := user.Cuisines[1:len(user.Cuisines)]
	if len(secondaryCuisines) > 0 {
		secondary = append(secondary, secondaryCuisines[0])
		if len(secondaryCuisines) > 1 {
			secondary = append(secondary, secondaryCuisines[1])
		}
	}
	return
}

func getPrimarySecondaryCost(user models.User) (primary models.CostTracking, secondary []models.CostTracking) {
	if len(user.CostBracket) == 0 {
		return
	}
	sort.Slice(user.CostBracket, func(i, j int) bool {
		return user.CostBracket[i].NoOfOrders > user.CostBracket[j].NoOfOrders
	})
	primary = user.CostBracket[0]
	secondaryCosts := user.CostBracket[1:len(user.CostBracket)]
	if len(secondaryCosts) > 0 {
		secondary = append(secondary, secondaryCosts[0])
		if len(secondaryCosts) > 1 {
			secondary = append(secondary, secondaryCosts[1])
		}
	}
	return
}
