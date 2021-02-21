const bcrypt = require('bcrypt');
const rounds = 12;

const hash = bcrypt.hashSync('secret', rounds);
console.log(hash);
