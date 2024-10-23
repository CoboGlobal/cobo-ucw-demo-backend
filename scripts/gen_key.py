from ed25519.keys import create_keypair

# Create a key pair
sk, vk = create_keypair()

# Print the private key in hexadecimal format.
print("private key:", sk.to_seed().hex())

# Print the public key in hexadecimal format.
print("public key:", vk.to_bytes().hex())
