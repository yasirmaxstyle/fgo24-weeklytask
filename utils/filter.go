package utils

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"
)

func (cli *CLI) filterMenu() {
	filter := cli.getFilterOptions()
	if filter == nil {
		return
	}

	var filteredItems []MenuItem
	for _, category := range cli.menu.MenuCategories {
		// Check if category is selected
		categorySelected := slices.Contains(filter.Categories, category.ID)

		if !categorySelected {
			continue
		}

		for _, item := range category.Items {
			if cli.matchesFilter(item, *filter) {
				filteredItems = append(filteredItems, item)
			}
		}
	}

	// sort filtered items by rating in descending order
	sort.Slice(filteredItems, func(i, j int) bool {
		return filteredItems[i].Rating > filteredItems[j].Rating
	})

	cli.clearScreen()
	cli.displayHeader()
	fmt.Println("FILTERED RESULTS")
	fmt.Printf("Found %d items\n\n", len(filteredItems))

	if len(filteredItems) == 0 {
		fmt.Println("No items match your filter criteria")
		cli.waitForEnter()
		return
	}

	filterCategory := MenuCategory{
		Name:  "Filter Results",
		Items: filteredItems,
	}

	if len(filteredItems) >= ItemsPerPage {
		cli.displayMenu(filterCategory)
	} else {
		for idx, item := range filteredItems {
			cli.displayMenuItem(item, true, idx)
		}
	}

	fmt.Println("\n0. Back to Main Menu")
	fmt.Print("\nSelect item to add to cart (or back): ")

	cli.scanner.Scan()
	choice, err := strconv.Atoi(cli.scanner.Text())
	if err != nil || choice < 1 {
		return
	}

	if choice == 0 {
		return
	}

	if choice <= len(filteredItems) {
		cli.addToCart(filteredItems[choice-1])
	}
}

// Get filter options using checkbox-style interface
func (cli *CLI) getFilterOptions() *Filter {
	cli.clearScreen()
	cli.displayHeader()

	fmt.Println("FILTER OPTIONS")
	fmt.Println("Select categories (enter numbers separated by commas, or press Enter for all):")

	for i, category := range cli.menu.MenuCategories {
		fmt.Printf("%d. %s\n", i+1, category.Name)
	}

	fmt.Print("\nCategories: ")
	cli.scanner.Scan()
	categoryInput := strings.TrimSpace(cli.scanner.Text())

	var selectedCategories []string
	if categoryInput != "" {
		parts := strings.Split(categoryInput, ",")
		for _, part := range parts {
			index, err := strconv.Atoi(strings.TrimSpace(part))
			if err == nil && index >= 1 && index <= len(cli.menu.MenuCategories) {
				selectedCategories = append(selectedCategories, cli.menu.MenuCategories[index-1].ID)
			}
		}
	} else {
		for _, category := range cli.menu.MenuCategories {
			selectedCategories = append(selectedCategories, category.ID)
		}
	}

	return &Filter{
		Categories: selectedCategories,
		Available:  true,
	}
}

// Check if item matches filter
func (cli *CLI) matchesFilter(item MenuItem, filter Filter) bool {
	if filter.Available && !item.Available {
		return false
	}

	return true
}
