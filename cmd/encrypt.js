const crypto = require("crypto");
const CryptoJS = require("crypto-js");

function generateDerivedKey(key, salt) {
  return crypto.pbkdf2Sync(key, salt, 1000, 16, "sha256").toString("hex");
}

function convertToAES(inputJSON, key) {
  const minifiedJSON = JSON.stringify(inputJSON);
  const base64 = Buffer.from(minifiedJSON).toString("base64");
  const salt = crypto.randomBytes(20).toString("hex");
  const derivedKey = generateDerivedKey(key, salt);
  const aesString = CryptoJS.AES.encrypt(base64, derivedKey).toString();
  const sha256key = crypto.createHash("sha256").update(key).digest("hex");

  return JSON.stringify({ derivedKey, salt, aesString, sha256key });
}


if (require.main === module) {
    const args = process.argv.slice(2);
    if (args.length < 2) {
      console.error("Usage: node encrypt.js '<JSON_DATA>' '<KEY>'");
      process.exit(1);
    }
  
    try {
      const jsonData = JSON.parse(args[0]);
      const key = args[1];
      console.log(convertToAES(jsonData, key));
    } catch (error) {
      console.error("Encryption Error:", error.message);
      process.exit(1);
    }
  }
  

module.exports = { convertToAES };
