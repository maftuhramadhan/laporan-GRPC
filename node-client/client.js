
const grpc = require("@grpc/grpc-js");
const protoLoader = require("@grpc/proto-loader");
const path = require("path");

const PROTO_PATH = path.join(__dirname, "..", "user.proto");

const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
    keepCase: false,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true
});

const protoDescriptor = grpc.loadPackageDefinition(packageDefinition);
const userPackage = protoDescriptor.user;

const client = new userPackage.UserService(
    "localhost:50051",
    grpc.credentials.createInsecure()
);

client.GetUserProfile({ userId: 1 }, (error, response) => {
    if (error) {
        console.error("Gagal mengambil data user:", error.message);
        client.close();
        process.exitCode = 1;
        return;
    }

    console.log("=== HASIL DATA USER ===");
    console.log("Nama  :", response.name);
    console.log("Email :", response.email);

    client.close();
});
