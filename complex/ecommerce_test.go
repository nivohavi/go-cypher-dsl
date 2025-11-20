package complex

import (
	"strings"
	"testing"

	"github.com/nivohavi/go-cypher-dsl/pkg/cypher"
	"github.com/nivohavi/go-cypher-dsl/pkg/cypher/ast"
)

// Test 12: Customer Purchase Analysis
func TestCustomerPurchaseAnalysis(t *testing.T) {
	customer := ast.Node("Customer").Named("customer")
	order := ast.Node("Order").Named("order")
	product := ast.Node("Product").Named("product")

	rel1 := customer.RelationshipTo(order, "PURCHASED")
	rel2 := order.RelationshipTo(product, "CONTAINS")

	stmt, err := cypher.Match(ast.Pattern(customer, rel1, order, rel2, product)).
		Returning(
			customer.Property("name"),
			cypher.As(cypher.Count(cypher.Distinct(order)), "totalOrders"),
			cypher.As(cypher.Collect(cypher.Distinct(product.Property("name"))), "productsPurchased"),
		).
		OrderBy(cypher.Desc(cypher.Var("totalOrders"))).
		Build()

	if err != nil {
		t.Fatalf("CustomerPurchaseAnalysis query Build() error = %v", err)
	}

	cypherStr := stmt.Cypher()
	if !strings.Contains(cypherStr, "PURCHASED") {
		t.Errorf("CustomerPurchaseAnalysis query should contain PURCHASED")
	}
	if !strings.Contains(cypherStr, "collect") {
		t.Errorf("CustomerPurchaseAnalysis query should contain collect")
	}
}

// Test 13: Product Category Analysis
func TestProductCategoryAnalysis(t *testing.T) {
	customer := ast.Node("Customer").Named("customer")
	order := ast.Node("Order").Named("order")
	product := ast.Node("Product").Named("product")
	category := ast.Node("Category").Named("category")

	rel1 := customer.RelationshipTo(order, "PLACED")
	rel2 := order.RelationshipTo(product, "CONTAINS")
	rel3 := product.RelationshipTo(category, "PART_OF")

	stmt, err := cypher.Match(ast.Pattern(customer, rel1, order, rel2, product, rel3, category)).
		Returning(
			customer.Property("name"),
			category.Property("name"),
			cypher.As(cypher.Count(cypher.Distinct(product)), "productCount"),
		).
		OrderBy(cypher.Desc(cypher.Var("productCount"))).
		Build()

	if err != nil {
		t.Fatalf("ProductCategoryAnalysis query Build() error = %v", err)
	}

	cypherStr := stmt.Cypher()
	if !strings.Contains(cypherStr, "PART_OF") {
		t.Errorf("ProductCategoryAnalysis query should contain PART_OF")
	}
}

// Test 14: Frequently Bought Together
func TestFrequentlyBoughtTogether(t *testing.T) {
	p1 := ast.Node("Product").Named("p1")
	p2 := ast.Node("Product").Named("p2")
	order := ast.Node("Order").Named("order")

	rel1 := order.RelationshipTo(p1, "CONTAINS")
	rel2 := order.RelationshipTo(p2, "CONTAINS")

	stmt, err := cypher.Match(ast.Pattern(order, rel1, p1)).
		Match(ast.Pattern(order, rel2, p2)).
		Where(
			cypher.Ne(cypher.Var("p1"), cypher.Var("p2")),
		).
		Returning(
			cypher.As(p1.Property("name"), "Product1"),
			cypher.As(p2.Property("name"), "Product2"),
			cypher.As(cypher.Count(cypher.Var("order")), "timesTogether"),
		).
		OrderBy(cypher.Desc(cypher.Var("timesTogether"))).
		Limit(10).
		Build()

	if err != nil {
		t.Fatalf("FrequentlyBoughtTogether query Build() error = %v", err)
	}

	cypherStr := stmt.Cypher()
	if !strings.Contains(cypherStr, "CONTAINS") {
		t.Errorf("FrequentlyBoughtTogether query should contain CONTAINS")
	}
}

