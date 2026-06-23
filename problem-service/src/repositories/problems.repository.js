const { pool } = require('../config/postgres');

const createProblemRecord = async (problemRecord) => {
  const query = `
    INSERT INTO problems (
      id,
      polygon_id,
      polygon_version,

      name,
      statement,
      input_statement,
      output_statement,
      examples,
      notes,

      memory_limit,
      time_limit,
      test_count,
      
      input_type,
      output_type,
      author_name,
      
      checker_language,
      checker_file_key,
      testcases_zip_key,
    )
      VALUES ($1::uuid, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
    RETURNING *
  `;

  const values = [
    problemRecord.id,
    problemRecord.polygonId,
    problemRecord.polygonVersion,

    problemRecord.name,
    problemRecord.statement,
    problemRecord.inputStatement,
    problemRecord.outputStatement,
    problemRecord.examples,
    problemRecord.notes,

    problemRecord.memoryLimit,
    problemRecord.timeLimit,
    problemRecord.testCount,

    problemRecord.inputType,
    problemRecord.outputType,
    problemRecord.authorName,
    
    problemRecord.checkerLanguage,
    problemRecord.checkerFileKey,
    problemRecord.testcasesZipKey,
  ];

  const { rows } = await pool.query(query, values);
  return rows[0];
};

const getProblemById = async (problemId) => {
  const { rows } = await pool.query('SELECT * FROM problems WHERE id = $1::uuid', [problemId]);
  return rows[0] ?? null;
};

const getProblemByName = async (name) => {
  const { rows } = await pool.query('SELECT * FROM problems WHERE name = $1', [name]);
  return rows[0] ?? null;
};

module.exports = {
  createProblemRecord,
  getProblemById,
  getProblemByName,
};