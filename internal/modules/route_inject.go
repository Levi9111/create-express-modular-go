package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	ImportMarker = "// --- INJECT IMPORTS HERE ---"
	RouteMarker  = "// --- INJECT ROUTES HERE ---"
)

// InjectRoute injects the import and route line into src/app/routes/index.ts.
func InjectRoute(projectRoot, importLine, routeLine string) error {
	indexPath := filepath.Join(projectRoot, "src/app/routes/index.ts")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return err
	}

	content := string(data)
	if !strings.Contains(content, ImportMarker) || !strings.Contains(content, RouteMarker) {
		return fmt.Errorf("route inject markers missing in src/app/routes/index.ts")
	}

	content = strings.Replace(content, ImportMarker, importLine+"\n"+ImportMarker, 1)
	content = strings.Replace(content, RouteMarker, routeLine+"\n  "+RouteMarker, 1)

	return os.WriteFile(indexPath, []byte(content), 0644)
}

// RemoveRoute unwires a module route and its import from src/app/routes/index.ts.
func RemoveRoute(projectRoot, moduleName string) error {
	indexPath := filepath.Join(projectRoot, "src/app/routes/index.ts")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return err
	}

	content := string(data)

	// Remove import line like: import { ProductRoutes } from '../modules/Product/product.route';
	importRegex := regexp.MustCompile(fmt.Sprintf(`(?m)^import\s*\{\s*%sRoutes\s*\}\s*from\s*['"].*['"];?\r?\n?`, moduleName))
	content = importRegex.ReplaceAllString(content, "")

	// Remove route entry like: { path: '/products', route: ProductRoutes },
	routeRegex := regexp.MustCompile(fmt.Sprintf(`(?m)^\s*\{\s*path:\s*['"][^'"]*['"],\s*route:\s*%sRoutes\s*\},?\r?\n?`, moduleName))
	content = routeRegex.ReplaceAllString(content, "")

	return os.WriteFile(indexPath, []byte(content), 0644)
}
