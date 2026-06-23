const fs = require('fs/promises');
const path = require('node:path');

const { exec } = require('child_process');
const execPromise = promisify(exec);
const MAX_TESTS_GENERATION_TIME = 5 * 60 * 1000; //5min

const generateTestCases = async (workDir) => {
  //ensure wine is installed
  try {
    await execPromise('wine --version');
  } catch (error) {
    throw new Error('Wine is required to generate testcases but it is not installed.');
  }

  //this is a very bad way to do it but giving permission is necessary
  const files = await fs.readdir(path.join(workDir, 'scripts'));
  await Promise.all(files.map(async (file) => {
    const fullPath = path.join(workDir, 'scripts', file);
    const fileStat = await fs.stat(fullPath);

    if (file.endsWith('.sh')) {
      // 0o755 gives read/write/execute to owner, and read/execute to others
      await fs.chmod(fullPath, 0o755);
    }
  }));

  //generates testcases
  try {
    const testcasesGenOutput = await execPromise('bash doall.sh', { cwd: workDir, timeout: MAX_TESTS_GENERATION_TIME });
  } catch (error) {
    throw new Error(`Testcase generation failed: ${error.message}, \nstderr: ${error.stderr}`);
  }
};

const validateTestCases = async (workDir, testSetName, testCount) => {
  const testsDir = path.join(workDir, testSetName);

  const entries = await fs.readdir(testsDir, { withFileTypes: true });

  if (entries.length !== testCount * 2) {
    return new Error(`Expected ${testCount * 2} test files, but found ${entries.length}.`);
  }

  for (const entry of entries) {
    const fullPath = path.join(testsDir, entry.name);

    if (!entry.isFile()) {
      throw new Error(`Only files are allowed in the test directory: ${fullPath}`);
    }
  }

  for (let i = 1; i <= testCount; i++) {
    if (!entries.some(e => e.name === `${i}`)) {
      throw new Error(`Missing input test file : ${i}`);
    }
    if (!entries.some(e => e.name === `${i}.a`)) {
      throw new Error(`Missing output test file : ${i}.a`);
    }
  }
};

module.exports = {
  generateTestCases,
  validateTestCases
}