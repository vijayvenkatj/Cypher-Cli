const crypto = require("crypto");
const CryptoJS = require("crypto-js");

function generateDerivedKey(key, salt) {
  return crypto.pbkdf2Sync(key, salt, 1000, 16, "sha256").toString("hex");
}

function convertFromAES(aesString, key, salt) {
  const derivedKey = generateDerivedKey(key, salt);
  const bytes = CryptoJS.AES.decrypt(aesString, derivedKey);

  if (!bytes || bytes.sigBytes <= 0) {
    console.error("Decryption failed.");
    process.exit(1);
  }

  try {
    
    const decryptedText = bytes.toString(CryptoJS.enc.Utf8);
    if (!decryptedText) {
      throw new Error("Decrypted text is empty.");
    }

    const decoded = Buffer.from(decryptedText, "base64").toString();
    return JSON.parse(decoded);
    
  } catch (error) {
    console.error(JSON.stringify({ error: "JSON Parse Error", details: error.message }));
    process.exit(1);
  }
}

// CLI Execution
if (require.main === module) {
  const [aesString, key, salt] = process.argv.slice(2);

  if (!aesString || !key || !salt) {
    console.error(JSON.stringify({ error: "Usage: node decrypt.js <aesString> <key> <salt>" }));
    process.exit(1);
  }

  const result = convertFromAES(aesString, key, salt);
  console.log(JSON.stringify(result)); // Always output pure JSON
}
