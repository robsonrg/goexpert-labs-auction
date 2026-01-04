const users = [
  { _id: "fd33527b-e74c-4235-8663-d34931a2936a", name: "David" },
  { _id: "1c3c9c61-002d-4566-a2a2-ae6d4f728c22", name: "Eve" }
];

const dbName = process.env.MONGODB_DB || "auctions";
const db = db.getSiblingDB(dbName);

users.forEach(u => {
  db.users.updateOne(
    { _id: u._id },
    { $setOnInsert: { name: u.name } },
    { upsert: true }
  );
});
