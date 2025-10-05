const fs = require('fs');
const path = require('path');

/**
 * Compares current PR test results with main branch results
 * @param {string} currentResultsPath - Path to current test results
 * @param {string} baselineResultsPath - Path to baseline (main branch) results
 * @returns {Object} Comparison summary
 */
function compareTestResults(currentResultsPath, baselineResultsPath) {
  try {
    const currentResults = JSON.parse(fs.readFileSync(currentResultsPath, 'utf8'));
    
    // If baseline doesn't exist, return current results only
    if (!fs.existsSync(baselineResultsPath)) {
      return {
        hasBaseline: false,
        current: extractTestSummary(currentResults),
        comparison: null,
        newFailures: [],
        fixedTests: [],
        newTests: [],
        removedTests: []
      };
    }
    
    const baselineResults = JSON.parse(fs.readFileSync(baselineResultsPath, 'utf8'));
    
    const currentSummary = extractTestSummary(currentResults);
    const baselineSummary = extractTestSummary(baselineResults);
    
    // Compare test results
    const comparison = {
      hasBaseline: true,
      current: currentSummary,
      baseline: baselineSummary,
      comparison: {
        totalChange: currentSummary.total - baselineSummary.total,
        passedChange: currentSummary.passed - baselineSummary.passed,
        failedChange: currentSummary.failed - baselineSummary.failed,
        flakyChange: currentSummary.flaky - baselineSummary.flaky,
        skippedChange: currentSummary.skipped - baselineSummary.skipped,
        durationChange: currentSummary.duration - baselineSummary.duration
      },
      newFailures: findNewFailures(currentResults, baselineResults),
      fixedTests: findFixedTests(currentResults, baselineResults),
      newTests: findNewTests(currentResults, baselineResults),
      removedTests: findRemovedTests(currentResults, baselineResults)
    };
    
    return comparison;
    
  } catch (error) {
    console.error('Error comparing test results:', error);
    return {
      hasBaseline: false,
      error: error.message,
      current: null,
      comparison: null
    };
  }
}

/**
 * Extracts test summary from Playwright results
 */
function extractTestSummary(results) {
  const stats = results.stats || {};
  return {
    total: (stats.expected || 0) + (stats.unexpected || 0) + (stats.flaky || 0) + (stats.skipped || 0),
    passed: stats.expected || 0,
    failed: stats.unexpected || 0,
    flaky: stats.flaky || 0,
    skipped: stats.skipped || 0,
    duration: stats.duration || 0
  };
}

/**
 * Finds tests that are now failing but were passing in baseline
 */
function findNewFailures(currentResults, baselineResults) {
  const currentTests = extractTestList(currentResults);
  const baselineTests = extractTestList(baselineResults);
  
  const newFailures = [];
  
  currentTests.forEach(test => {
    if (test.status === 'failed') {
      const baselineTest = baselineTests.find(bt => 
        bt.title === test.title && bt.file === test.file
      );
      
      if (baselineTest && baselineTest.status === 'passed') {
        newFailures.push({
          title: test.title,
          file: test.file,
          project: test.project,
          error: test.error
        });
      }
    }
  });
  
  return newFailures;
}

/**
 * Finds tests that were failing but are now passing
 */
function findFixedTests(currentResults, baselineResults) {
  const currentTests = extractTestList(currentResults);
  const baselineTests = extractTestList(baselineResults);
  
  const fixedTests = [];
  
  currentTests.forEach(test => {
    if (test.status === 'passed') {
      const baselineTest = baselineTests.find(bt => 
        bt.title === test.title && bt.file === test.file
      );
      
      if (baselineTest && baselineTest.status === 'failed') {
        fixedTests.push({
          title: test.title,
          file: test.file,
          project: test.project
        });
      }
    }
  });
  
  return fixedTests;
}

/**
 * Finds tests that exist in current but not in baseline
 */
function findNewTests(currentResults, baselineResults) {
  const currentTests = extractTestList(currentResults);
  const baselineTests = extractTestList(baselineResults);
  
  const newTests = [];
  
  currentTests.forEach(test => {
    const baselineTest = baselineTests.find(bt => 
      bt.title === test.title && bt.file === test.file
    );
    
    if (!baselineTest) {
      newTests.push({
        title: test.title,
        file: test.file,
        project: test.project,
        status: test.status
      });
    }
  });
  
  return newTests;
}

/**
 * Finds tests that exist in baseline but not in current
 */
function findRemovedTests(currentResults, baselineResults) {
  const currentTests = extractTestList(currentResults);
  const baselineTests = extractTestList(baselineResults);
  
  const removedTests = [];
  
  baselineTests.forEach(test => {
    const currentTest = currentTests.find(ct => 
      ct.title === test.title && ct.file === test.file
    );
    
    if (!currentTest) {
      removedTests.push({
        title: test.title,
        file: test.file,
        project: test.project,
        status: test.status
      });
    }
  });
  
  return removedTests;
}

/**
 * Extracts a flat list of tests from Playwright results
 */
function extractTestList(results) {
  const tests = [];
  
  function traverseSuites(suites, parentFile = '') {
    suites.forEach(suite => {
      const file = suite.file || parentFile;
      
      if (suite.specs) {
        suite.specs.forEach(spec => {
          spec.tests?.forEach(test => {
            const results = test.results || [];
            const status = results.length > 0 ? results[results.length - 1].status : 'unknown';
            const error = results.find(r => r.status === 'failed')?.error?.message;
            
            tests.push({
              title: test.title,
              file: file,
              project: test.projectName || 'unknown',
              status: status,
              error: error
            });
          });
        });
      }
      
      if (suite.suites) {
        traverseSuites(suite.suites, file);
      }
    });
  }
  
  if (results.suites) {
    traverseSuites(results.suites);
  }
  
  return tests;
}

/**
 * Formats comparison results into markdown
 */
function formatComparison(comparison) {
  if (!comparison.hasBaseline) {
    return `### 📊 Test Comparison\n\n*No baseline results available for comparison*\n\n`;
  }
  
  let markdown = `### 📊 Test Comparison vs Main Branch\n\n`;
  
  // Summary comparison table
  markdown += `| Metric | Current | Baseline | Change |\n`;
  markdown += `|--------|---------|----------|--------|\n`;
  
  const { current, baseline, comparison: comp } = comparison;
  
  markdown += `| **Total Tests** | ${current.total} | ${baseline.total} | ${formatChange(comp.totalChange)} |\n`;
  markdown += `| **✅ Passed** | ${current.passed} | ${baseline.passed} | ${formatChange(comp.passedChange)} |\n`;
  markdown += `| **❌ Failed** | ${current.failed} | ${baseline.failed} | ${formatChange(comp.failedChange)} |\n`;
  markdown += `| **⚠️ Flaky** | ${current.flaky} | ${baseline.flaky} | ${formatChange(comp.flakyChange)} |\n`;
  markdown += `| **⏭️ Skipped** | ${current.skipped} | ${baseline.skipped} | ${formatChange(comp.skippedChange)} |\n`;
  
  const durationChangeMin = Math.round(comp.durationChange / 1000 / 60 * 100) / 100;
  markdown += `| **⏱️ Duration** | ${Math.round(current.duration / 1000 / 60 * 100) / 100}m | ${Math.round(baseline.duration / 1000 / 60 * 100) / 100}m | ${formatChange(durationChangeMin, 'm')} |\n\n`;
  
  // New failures
  if (comparison.newFailures.length > 0) {
    markdown += `#### 🚨 New Failures (${comparison.newFailures.length})\n\n`;
    comparison.newFailures.forEach(test => {
      markdown += `- **${test.title}** in \`${test.file}\`\n`;
    });
    markdown += `\n`;
  }
  
  // Fixed tests
  if (comparison.fixedTests.length > 0) {
    markdown += `#### ✅ Fixed Tests (${comparison.fixedTests.length})\n\n`;
    comparison.fixedTests.forEach(test => {
      markdown += `- **${test.title}** in \`${test.file}\`\n`;
    });
    markdown += `\n`;
  }
  
  // New tests
  if (comparison.newTests.length > 0) {
    markdown += `#### 🆕 New Tests (${comparison.newTests.length})\n\n`;
    comparison.newTests.forEach(test => {
      const statusIcon = test.status === 'passed' ? '✅' : test.status === 'failed' ? '❌' : '⚠️';
      markdown += `- ${statusIcon} **${test.title}** in \`${test.file}\`\n`;
    });
    markdown += `\n`;
  }
  
  // Removed tests
  if (comparison.removedTests.length > 0) {
    markdown += `#### 🗑️ Removed Tests (${comparison.removedTests.length})\n\n`;
    comparison.removedTests.forEach(test => {
      markdown += `- **${test.title}** in \`${test.file}\`\n`;
    });
    markdown += `\n`;
  }
  
  return markdown;
}

/**
 * Formats a numeric change with appropriate emoji and sign
 */
function formatChange(change, unit = '') {
  if (change === 0) return `0${unit}`;
  
  const sign = change > 0 ? '+' : '';
  const emoji = change > 0 ? '📈' : '📉';
  
  return `${sign}${change}${unit} ${emoji}`;
}

// CLI usage
if (require.main === module) {
  const currentPath = process.argv[2];
  const baselinePath = process.argv[3];
  
  if (!currentPath) {
    console.error('Usage: node compare-results.js <current-results.json> [baseline-results.json]');
    process.exit(1);
  }
  
  const comparison = compareTestResults(currentPath, baselinePath);
  const markdown = formatComparison(comparison);
  console.log(markdown);
}

module.exports = { compareTestResults, formatComparison };