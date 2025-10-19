#!/usr/bin/env ts-node

/**
 * MCP Playwright Server
 * 
 * This server exposes Playwright E2E testing capabilities through the Model Context Protocol (MCP).
 * It allows AI assistants to run, debug, and manage Playwright tests programmatically.
 */

import { Server } from '@modelcontextprotocol/sdk/server/index.js';
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js';
import {
  CallToolRequestSchema,
  ListToolsRequestSchema,
} from '@modelcontextprotocol/sdk/types.js';
import { exec } from 'child_process';
import { promisify } from 'util';
import * as path from 'path';
import * as fs from 'fs/promises';
import { fileURLToPath } from 'url';
import { dirname } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const execAsync = promisify(exec);

/**
 * MCP Server for Playwright Test Management
 */
class PlaywrightMCPServer {
  private server: Server;
  private projectRoot: string;

  constructor() {
    this.projectRoot = path.resolve(__dirname, '../..');
    
    this.server = new Server(
      {
        name: 'playwright-e2e-server',
        version: '1.0.0',
      },
      {
        capabilities: {
          tools: {},
        },
      }
    );

    this.setupHandlers();
  }

  private setupHandlers(): void {
    // List available tools
    this.server.setRequestHandler(ListToolsRequestSchema, async () => {
      return {
        tools: this.getTools(),
      };
    });

    // Handle tool execution
    this.server.setRequestHandler(CallToolRequestSchema, async (request) => {
      const { name, arguments: args } = request.params;

      try {
        switch (name) {
          case 'run_playwright_tests':
            return await this.runTests(args);
          
          case 'run_playwright_test_file':
            return await this.runTestFile(args);
          
          case 'list_playwright_tests':
            return await this.listTests();
          
          case 'get_test_report':
            return await this.getTestReport(args);
          
          case 'debug_playwright_test':
            return await this.debugTest(args);
          
          case 'run_test_with_ui':
            return await this.runTestWithUI(args);
          
          default:
            throw new Error(`Unknown tool: ${name}`);
        }
      } catch (error) {
        return {
          content: [
            {
              type: 'text',
              text: `Error executing ${name}: ${error instanceof Error ? error.message : String(error)}`,
            },
          ],
          isError: true,
        };
      }
    });
  }

  private getTools(): any[] {
    return [
      {
        name: 'run_playwright_tests',
        description: 'Run all Playwright E2E tests or filtered by pattern',
        inputSchema: {
          type: 'object',
          properties: {
            grep: {
              type: 'string',
              description: 'Filter tests by pattern (e.g., "subscriptions" or "happy-path")',
            },
            headed: {
              type: 'boolean',
              description: 'Run tests in headed mode (visible browser)',
              default: false,
            },
            project: {
              type: 'string',
              description: 'Run tests for specific project (e.g., "chromium", "setup")',
            },
          },
        },
      },
      {
        name: 'run_playwright_test_file',
        description: 'Run a specific Playwright test file',
        inputSchema: {
          type: 'object',
          properties: {
            testFile: {
              type: 'string',
              description: 'Path to test file relative to tests/e2e (e.g., "auth/happy-path.spec.ts")',
            },
            headed: {
              type: 'boolean',
              description: 'Run test in headed mode (visible browser)',
              default: false,
            },
          },
          required: ['testFile'],
        },
      },
      {
        name: 'list_playwright_tests',
        description: 'List all available Playwright test files with their descriptions',
        inputSchema: {
          type: 'object',
          properties: {},
        },
      },
      {
        name: 'get_test_report',
        description: 'Get the latest test execution report',
        inputSchema: {
          type: 'object',
          properties: {
            format: {
              type: 'string',
              enum: ['json', 'junit', 'html'],
              description: 'Report format to retrieve',
              default: 'json',
            },
          },
        },
      },
      {
        name: 'debug_playwright_test',
        description: 'Run a specific test in debug mode with inspector',
        inputSchema: {
          type: 'object',
          properties: {
            testFile: {
              type: 'string',
              description: 'Path to test file relative to tests/e2e',
            },
            testName: {
              type: 'string',
              description: 'Specific test name to debug (optional)',
            },
          },
          required: ['testFile'],
        },
      },
      {
        name: 'run_test_with_ui',
        description: 'Open Playwright UI mode for interactive test running',
        inputSchema: {
          type: 'object',
          properties: {
            testFile: {
              type: 'string',
              description: 'Optional specific test file to focus on',
            },
          },
        },
      },
    ];
  }

  private async runTests(args: any): Promise<any> {
    const { grep, headed, project } = args;
    
    let command = 'npm run test:e2e --';
    
    if (grep) {
      command += ` --grep "${grep}"`;
    }
    
    if (headed) {
      command += ' --headed';
    }
    
    if (project) {
      command += ` --project="${project}"`;
    }

    try {
      const { stdout, stderr } = await execAsync(command, {
        cwd: this.projectRoot,
        maxBuffer: 10 * 1024 * 1024, // 10MB buffer
      });

      return {
        content: [
          {
            type: 'text',
            text: `Test execution completed:\n\n${stdout}\n${stderr ? `\nWarnings/Errors:\n${stderr}` : ''}`,
          },
        ],
      };
    } catch (error: any) {
      return {
        content: [
          {
            type: 'text',
            text: `Test execution failed:\n\nCommand: ${command}\n\nOutput:\n${error.stdout || ''}\n\nError:\n${error.stderr || error.message}`,
          },
        ],
        isError: true,
      };
    }
  }

  private async runTestFile(args: any): Promise<any> {
    const { testFile, headed } = args;
    
    const testPath = path.join('tests/e2e', testFile);
    let command = `npm run test:e2e -- ${testPath}`;
    
    if (headed) {
      command += ' --headed';
    }

    try {
      const { stdout, stderr } = await execAsync(command, {
        cwd: this.projectRoot,
        maxBuffer: 10 * 1024 * 1024,
      });

      return {
        content: [
          {
            type: 'text',
            text: `Test file execution completed:\n\n${stdout}\n${stderr ? `\nWarnings:\n${stderr}` : ''}`,
          },
        ],
      };
    } catch (error: any) {
      return {
        content: [
          {
            type: 'text',
            text: `Test file execution failed:\n\nFile: ${testFile}\n\nOutput:\n${error.stdout || ''}\n\nError:\n${error.stderr || error.message}`,
          },
        ],
        isError: true,
      };
    }
  }

  private async listTests(): Promise<any> {
    const testsDir = path.join(this.projectRoot, 'tests/e2e');
    
    try {
      const testFiles = await this.findTestFiles(testsDir);
      const testList = testFiles.map(file => {
        const relativePath = path.relative(testsDir, file);
        return `- ${relativePath}`;
      }).join('\n');

      return {
        content: [
          {
            type: 'text',
            text: `Available Playwright E2E Tests:\n\n${testList}\n\nTotal: ${testFiles.length} test files`,
          },
        ],
      };
    } catch (error) {
      return {
        content: [
          {
            type: 'text',
            text: `Failed to list tests: ${error instanceof Error ? error.message : String(error)}`,
          },
        ],
        isError: true,
      };
    }
  }

  private async findTestFiles(dir: string): Promise<string[]> {
    const files: string[] = [];
    const entries = await fs.readdir(dir, { withFileTypes: true });

    for (const entry of entries) {
      const fullPath = path.join(dir, entry.name);
      
      if (entry.isDirectory()) {
        files.push(...await this.findTestFiles(fullPath));
      } else if (entry.isFile() && entry.name.endsWith('.spec.ts')) {
        files.push(fullPath);
      }
    }

    return files;
  }

  private async getTestReport(args: any): Promise<any> {
    const { format = 'json' } = args;
    
    const reportPaths: Record<string, string> = {
      json: 'test-results/test-results.json',
      junit: 'test-results/junit-report.xml',
      html: 'test-results/html-report/index.html',
    };

    const reportPath = path.join(this.projectRoot, reportPaths[format]);

    try {
      const reportContent = await fs.readFile(reportPath, 'utf-8');
      
      return {
        content: [
          {
            type: 'text',
            text: `Test Report (${format}):\n\n${reportContent}`,
          },
        ],
      };
    } catch (error) {
      return {
        content: [
          {
            type: 'text',
            text: `Failed to read test report: ${error instanceof Error ? error.message : String(error)}\n\nMake sure tests have been run at least once.`,
          },
        ],
        isError: true,
      };
    }
  }

  private async debugTest(args: any): Promise<any> {
    const { testFile, testName } = args;
    
    const testPath = path.join('tests/e2e', testFile);
    let command = `npm run test:e2e:debug -- ${testPath}`;
    
    if (testName) {
      command += ` --grep "${testName}"`;
    }

    return {
      content: [
        {
          type: 'text',
          text: `Starting debug mode for test:\n\nFile: ${testFile}${testName ? `\nTest: ${testName}` : ''}\n\nCommand: ${command}\n\nNote: This will open the Playwright Inspector. The server will wait for the debug session to complete.`,
        },
      ],
    };
  }

  private async runTestWithUI(args: any): Promise<any> {
    const { testFile } = args;
    
    let command = 'npm run test:e2e:ui';
    
    if (testFile) {
      command += ` -- tests/e2e/${testFile}`;
    }

    return {
      content: [
        {
          type: 'text',
          text: `Opening Playwright UI mode...\n\nCommand: ${command}\n\nThis will open an interactive UI for running and debugging tests.${testFile ? `\n\nFocused on: ${testFile}` : ''}`,
        },
      ],
    };
  }

  async run(): Promise<void> {
    const transport = new StdioServerTransport();
    await this.server.connect(transport);
    
    console.error('Playwright MCP Server running on stdio');
  }
}

// Start the server
const server = new PlaywrightMCPServer();
server.run().catch((error) => {
  console.error('Fatal error running server:', error);
  process.exit(1);
});
