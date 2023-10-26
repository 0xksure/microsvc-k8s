import * as anchor from "@coral-xyz/anchor";
import { Identity } from "../target/types/identity";
import * as identity from "../sdk-ts/src/index"
import { sendAndConfirmTransaction } from "../sdk-ts/src/utils";
import { assert, config, expect, use } from 'chai';
import * as chaiAsPromised from "chai-as-promised"
import { bs58 } from "@coral-xyz/anchor/dist/cjs/utils/bytes";
use(chaiAsPromised.default)
/**
 * topUpAccount is a helper function to top up an account with SOL
 * @param connection 
 * @param wallet 
 * @returns 
 */
const topUpAccount = async (connection: anchor.web3.Connection, wallet: anchor.Wallet) => {
  const latestBlockhash = await connection.getLatestBlockhash();

  const fromAirdropSig = await connection.requestAirdrop(
    wallet.publicKey,
    10 * anchor.web3.LAMPORTS_PER_SOL
  );
  return await connection.confirmTransaction({
    signature: fromAirdropSig,
    ...latestBlockhash
  });
}

describe("identity", () => {
  // Configure the client to use the local cluster.

  const wallet = new anchor.Wallet(anchor.web3.Keypair.generate());
  const user = new anchor.Wallet(anchor.web3.Keypair.generate());

  const program = anchor.workspace.Identity as anchor.Program<Identity>;
  const identitySdk = new identity.IdentitySdk(wallet.publicKey, program.provider.connection);
  const localAnchorProvider = anchor.AnchorProvider.env();
  const provider = new anchor.AnchorProvider(
    localAnchorProvider.connection,
    wallet,
    localAnchorProvider.opts
  );
  anchor.setProvider(provider);

  before(async () => {
    await topUpAccount(program.provider.connection, wallet);
    await topUpAccount(program.provider.connection, user);

    // initialize protocol

    const initializeIdentity = await identitySdk.initializeProtocol()
    await sendAndConfirmTransaction(
      program.provider.connection,
      initializeIdentity.vtx,
      [wallet.payer]

    )

    const identityProgramAccount = await program.account.identityProgram.fetch(initializeIdentity.identityProgramPDA[0])
    expect(identityProgramAccount.protocolOwner.toString()).to.equal(wallet.payer.publicKey.toString())


  })

  it("Create identity -> Should succeed ", async () => {
    const createIdentityParams = {
      social: "github",
      userId: 444322,
      username: "test",
      identityOwner: user.publicKey,
      protocolOwner: wallet.publicKey
    }
    const createIdentity = await identitySdk.createIdentity(createIdentityParams)
    await sendAndConfirmTransaction(
      program.provider.connection,
      createIdentity.vtx,
      [user.payer, wallet.payer],

    )

    const identity = await program.account.identity.fetch(createIdentity.identityPDA[0])
    expect(identity.social.github).to.not.be.undefined
    expect(identity.userId).to.equal(createIdentityParams.userId)
  })

  it("Try to recreate the protocol -> should fail", async () => {
    const initializeIdentity = await identitySdk.initializeProtocol()
    await expect(sendAndConfirmTransaction(
      program.provider.connection,
      initializeIdentity.vtx,
      [wallet.payer]

    )).to.be.rejectedWith(Error)
  })

  it("Try to create the same identity twice -> should fail", async () => {
    const createIdentityParams = {
      social: "github",
      userId: 444323,
      username: "test",
      identityOwner: user.publicKey,
      protocolOwner: wallet.publicKey
    }
    const createIdentity = await identitySdk.createIdentity(createIdentityParams)
    await sendAndConfirmTransaction(
      program.provider.connection,
      createIdentity.vtx,
      [user.payer, wallet.payer],
    )
    await expect(sendAndConfirmTransaction(
      program.provider.connection,
      createIdentity.vtx,
      [user.payer, wallet.payer],

    )).to.be.rejectedWith(Error)
  })

  it("Try to create an identity with a different protocol owner -> should fail", async () => {
    const notProtocolOwner = new anchor.Wallet(anchor.web3.Keypair.generate());
    const createIdentityParams = {
      social: "github",
      userId: 444324,
      username: "test",
      identityOwner: user.publicKey,
      protocolOwner: notProtocolOwner.publicKey
    }
    const createIdentity = await identitySdk.createIdentity(createIdentityParams)
    await expect(sendAndConfirmTransaction(
      program.provider.connection,
      createIdentity.vtx,
      [user.payer, notProtocolOwner.payer],

    )).to.be.rejectedWith(Error)
  })

  it("Try to create an identity with as only the identity owner -> should fail", async () => {
    const notIdentityOwner = new anchor.Wallet(anchor.web3.Keypair.generate());
    const createIdentityParams = {
      social: "github",
      userId: 444325,
      username: "test",
      identityOwner: user.publicKey,
      protocolOwner: user.publicKey
    }
    const createIdentity = await identitySdk.createIdentity(createIdentityParams)
    await expect(sendAndConfirmTransaction(
      program.provider.connection,
      createIdentity.vtx,
      [notIdentityOwner.payer],

    )).to.be.rejectedWith(Error)
  })

  it("Create identity and update username as identity owner-> should succeed", async () => {
    const createIdentityParams = {
      social: "github",
      userId: 444326,
      username: "test",
      identityOwner: user.publicKey,
      protocolOwner: wallet.publicKey
    }
    const createIdentity = await identitySdk.createIdentity(createIdentityParams)
    await sendAndConfirmTransaction(
      program.provider.connection,
      createIdentity.vtx,
      [user.payer, wallet.payer],

    )
    const firstCreatedIdentity = await program.account.identity.fetch(createIdentity.identityPDA[0])
    expect(firstCreatedIdentity.social.github).to.not.be.undefined
    expect(firstCreatedIdentity.userId).to.equal(createIdentityParams.userId)
    expect(firstCreatedIdentity.username.toString()).to.equal(createIdentityParams.username)


    const newUsername = "test2"
    const updateUsername = await identitySdk.updateUsername({
      username: newUsername,
      identityOwner: user.publicKey,
      social: "github",
      userId: 444326
    })
    await sendAndConfirmTransaction(
      program.provider.connection,
      updateUsername.vtx,
      [user.payer],

    )

    const identity = await program.account.identity.fetch(createIdentity.identityPDA[0])
    expect(identity.social.github).to.not.be.undefined
    expect(identity.userId).to.equal(createIdentityParams.userId)
    expect(identity.username.toString()).to.equal(newUsername)
  })

  it("Create identity and try to change the userId -> should fail", async () => {
    const createIdentityParams = {
      social: "github",
      userId: 444327,
      username: "test",
      identityOwner: user.publicKey,
      protocolOwner: wallet.publicKey
    }
    const createIdentity = await identitySdk.createIdentity(createIdentityParams)
    await sendAndConfirmTransaction(
      program.provider.connection,
      createIdentity.vtx,
      [user.payer, wallet.payer],

    )
    const firstCreatedIdentity = await program.account.identity.fetch(createIdentity.identityPDA[0])
    expect(firstCreatedIdentity.social.github).to.not.be.undefined
    expect(firstCreatedIdentity.userId).to.equal(createIdentityParams.userId)
    expect(firstCreatedIdentity.username.toString()).to.equal(createIdentityParams.username)

    // update userName with new userID -> should fail 
    const newUsername = "test2"
    const updateUsername = await identitySdk.updateUsername({
      username: newUsername,
      identityOwner: user.publicKey,
      social: "github",
      userId: 444328
    })
    await expect(sendAndConfirmTransaction(
      program.provider.connection,
      updateUsername.vtx,
      [user.payer],

    )).to.be.rejectedWith(Error)
  })

  it("Try to transfer ownership of the identity to another address -> should succeed", async () => {
    const createIdentityParams = {
      social: "github",
      userId: 444329,
      username: "test",
      identityOwner: user.publicKey,
      protocolOwner: wallet.publicKey
    }
    const createIdentity = await identitySdk.createIdentity(createIdentityParams)
    await sendAndConfirmTransaction(
      program.provider.connection,
      createIdentity.vtx,
      [user.payer, wallet.payer],

    )
    const firstCreatedIdentity = await program.account.identity.fetch(createIdentity.identityPDA[0])
    expect(firstCreatedIdentity.social.github).to.not.be.undefined
    expect(firstCreatedIdentity.userId).to.equal(createIdentityParams.userId)
    expect(firstCreatedIdentity.username.toString()).to.equal(createIdentityParams.username)
    expect(firstCreatedIdentity.address.toString()).to.equal(user.publicKey.toString())

    // transfer ownership of the identity to another address
    const newOwner = new anchor.Wallet(anchor.web3.Keypair.generate());
    const transferOwnership = await identitySdk.transferOwnership({
      identityOwner: user.publicKey,
      newIdentityOwner: newOwner.publicKey,
      social: "github",
      userId: 444329
    })
    await sendAndConfirmTransaction(
      program.provider.connection,
      transferOwnership.vtx,
      [user.payer, newOwner.payer],

    )

    const identity = await program.account.identity.fetch(createIdentity.identityPDA[0])
    expect(identity.social.github).to.not.be.undefined
    expect(identity.userId).to.equal(createIdentityParams.userId)
    expect(identity.address.toString()).to.equal(newOwner.publicKey.toString())
  })

  it.only("create identity and fetch it using the username ", async () => {
    const username = "partyOn"
    const createIdentityParams = {
      social: "github",
      userId: 444330,
      username,
      identityOwner: user.publicKey,
      protocolOwner: wallet.publicKey
    }
    const createIdentity = await identitySdk.createIdentity(createIdentityParams)
    await sendAndConfirmTransaction(
      program.provider.connection,
      createIdentity.vtx,
      [user.payer, wallet.payer],

    )
    const firstCreatedIdentity = await program.account.identity.fetch(createIdentity.identityPDA[0])
    // strip trailing zeros from username array
    expect(firstCreatedIdentity.social.toString()).to.not.be.undefined
    expect(firstCreatedIdentity.userId).to.equal(createIdentityParams.userId)
    expect(identity.convertStringOfSizeToString(firstCreatedIdentity.username)).to.equal(createIdentityParams.username)
    expect(firstCreatedIdentity.address.toString()).to.equal(user.publicKey.toString())



    const matchingUnAccounts = await identitySdk.getIdentityFromUsername({
      username
    })
    console.log("matchingUnAccounts: ", matchingUnAccounts)
    expect(matchingUnAccounts.length).to.equal(1)
    expect(identity.convertStringOfSizeToString(matchingUnAccounts[0].account.username)).to.equal(createIdentityParams.username)
  })

});
