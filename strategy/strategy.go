package strategy

import (
	"recommendations/models"
	"reflect"
	"sort"
)

func ContainsCostTracking(cTracking []models.CostTracking, costBracketType int) bool {
	for _, v := range cTracking {
		if v.Type == costBracketType {
			return true
		}
	}
	return false
}

func ConntainsCuisine(cuisines []models.CuisineTracking, cuisineType string) bool {
	for _, v := range cuisines {
		if v.CuisineType == cuisineType {
			return true
		}
	}
	return false
}

type strategies struct {
}

func (s strategies) RecomendedPrimaryCousineCost(user models.User, restaurant models.Restaurant) bool {
	return restaurant.IsRecommended && restaurant.Cuisine == user.PrimaryCuisine.CuisineType && restaurant.CostBracket == user.PrimaryCostBracket.Type
}

func (s strategies) PrimaryCuisineSecondaryCost(user models.User, restaurant models.Restaurant) bool {
	return restaurant.Cuisine == user.PrimaryCuisine.CuisineType && ContainsCostTracking(user.SecondaryCostBracket, restaurant.CostBracket)
}

func (s strategies) SecondaryCuisinePrimaryCost(user models.User, restaurant models.Restaurant) bool {
	return ConntainsCuisine(user.SecondaryCuisine, restaurant.Cuisine) && user.PrimaryCostBracket.Type == restaurant.CostBracket
}
func (s strategies) PrimaryCuisineCostRatingGTE4(user models.User, restaurant models.Restaurant) bool {
	return restaurant.Cuisine == user.PrimaryCuisine.CuisineType && restaurant.CostBracket == user.PrimaryCostBracket.Type && restaurant.Rating >= 4
}

func (s strategies) PrimaryCuisineSecondaryCostRatingnGTE4_5(user models.User, restaurant models.Restaurant) bool {
	return restaurant.Cuisine == user.PrimaryCuisine.CuisineType && ContainsCostTracking(user.SecondaryCostBracket, restaurant.CostBracket) && restaurant.Rating >= 4.5
}

func (s strategies) SecondaryCuisinePrimaryCostRatingGTE4_5(user models.User, restaurant models.Restaurant) bool {
	return ConntainsCuisine(user.SecondaryCuisine, restaurant.Cuisine) && restaurant.CostBracket == user.PrimaryCostBracket.Type && restaurant.Rating >= 4.5
}
func (s strategies) PrimaryCuisineCostRatingLT4(user models.User, restaurant models.Restaurant) bool {
	return restaurant.Cuisine == user.PrimaryCuisine.CuisineType && restaurant.CostBracket == user.PrimaryCostBracket.Type && restaurant.Rating > 4
}

func (s strategies) PrimaryCuisineSecondaryCostRatingLT4_5(user models.User, restaurant models.Restaurant) bool {
	return restaurant.Cuisine == user.PrimaryCuisine.CuisineType && ContainsCostTracking(user.SecondaryCostBracket, restaurant.CostBracket) && restaurant.Rating < 4.5
}

func (s strategies) SeconndaryCuisinePrimaryCostRatingLT4_5(user models.User, restaurant models.Restaurant) bool {
	return ConntainsCuisine(user.SecondaryCuisine, restaurant.Cuisine) && user.PrimaryCostBracket.Type == restaurant.CostBracket && restaurant.Rating < 4.5
}

func (s strategies) AnyCuisineCost(user models.User, restaurant models.Restaurant) bool {
	return true
}

func (s strategies) Top4NewlyCreatedByRating(user models.User, restaurant models.Restaurant) bool {
	return false
}

var StrategyList = []string{
	"RecomendedPrimaryCousineCost",
	"PrimaryCuisineSecondaryCost",
	"SecondaryCuisinePrimaryCost",
	"PrimaryCuisineCostRatingGTE4",
	"PrimaryCuisineSecondaryCostRatingnGTE4_5",
	"SecondaryCuisinePrimaryCostRatingGTE4_5",
	"Top4NewlyCreatedByRating",
	"PrimaryCuisineCostRatingLT4",
	"PrimaryCuisineSecondaryCostRatingLT4_5",
	"SeconndaryCuisinePrimaryCostRatingLT4_5",
	"AnyCuisineCost",
}

func GetRestaurantsBasedOnStrategies(user models.User, availableRestaurants []models.Restaurant) map[int][]models.Restaurant {
	orderResults := map[int][]models.Restaurant{}
	for i := 0; i < len(StrategyList); i++ {
		orderResults[i] = []models.Restaurant{}
	}
	for _, restaurant := range availableRestaurants {
		for i, v := range StrategyList {
			params := []reflect.Value{reflect.ValueOf(user), reflect.ValueOf(restaurant)}
			funcName := reflect.ValueOf(strategies{}).MethodByName(v)
			returnVal := funcName.Call(params)

			if len(returnVal) > 0 && returnVal[0].Interface().(bool) {
				orderResults[i] = append(orderResults[i], restaurant)
				break
			}
		}
	}
	sort.Slice(availableRestaurants, func(i, j int) bool {
		return availableRestaurants[i].OnboardedTime.After(availableRestaurants[j].OnboardedTime)
	})
	newlyOnboarded4 := availableRestaurants
	if len(availableRestaurants) > 4 {
		newlyOnboarded4 = availableRestaurants[0:4]
	}
	sort.Slice(newlyOnboarded4, func(i, j int) bool {
		return newlyOnboarded4[i].Rating > newlyOnboarded4[j].Rating
	})
	orderResults[6] = append(orderResults[6], newlyOnboarded4...)
	return orderResults
}
