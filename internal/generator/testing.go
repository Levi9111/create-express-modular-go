package generator

import (
	"os"
	"path/filepath"
)

func scaffoldTesting(opts ProjectOptions) error {
	vitestConfig := `import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    environment: 'node',
    include: ['src/**/*.test.ts'],
  },
});
`

	healthTest := `import { describe, it, expect } from 'vitest';
import request from 'supertest';
import app from '../app';

describe('Health Check', () => {
  it('GET /health returns 200 ok', async () => {
    const res = await request(app).get('/health');
    expect(res.status).toBe(200);
    expect(res.body.status).toBe('ok');
  });
});
`

	if err := writeFile(filepath.Join(opts.ProjectPath, "vitest.config.ts"), vitestConfig); err != nil {
		return err
	}
	testDir := filepath.Join(opts.ProjectPath, "src", "__tests__")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		return err
	}
	return writeFile(filepath.Join(testDir, "health.test.ts"), healthTest)
}

