use anchor_lang::prelude::*;
use std::mem::size_of;

declare_id!("3Nt1tyTJ6VBf4APaPPWixXFJr6DtfGvvTwHY1aGYT4Ws");

const IDENTITY_SEED: &str = "identity";

#[error_code]
#[derive(PartialEq)]
pub enum ErrorCodes {
    #[msg("The protocol owner is not the owner of the account")]
    ProtocolOwnerNotOwner,

    #[msg("The signer is not the owner of the account")]
    SignerNotOwner,

    #[msg("The username is too long. Should be max 32bytes")]
    UsernameTooLong,
}

#[program]
pub mod identity {
    use super::*;

    /// Initialize the identity program
    /// Sets the protocol owner that is allowed to create identities
    pub fn initialize(ctx: Context<Initialize>) -> Result<()> {
        let identity_program = &mut ctx.accounts.identity_program;
        identity_program.protocol_owner = *ctx.accounts.protocol_owner.key;
        identity_program.bump = *ctx.bumps.get("identity_program").unwrap();
        Ok(())
    }

    /// Create a new web2 <-> web3 identity
    pub fn create_identity(
        ctx: Context<CreateIdentity>,
        social: String,
        username: String,
        user_id: u64,
    ) -> Result<()> {
        let identity = &mut ctx.accounts.identity;
        let bump = ctx.bumps.get("identity").unwrap();
        identity.init(
            *ctx.accounts.account_holder.key,
            social,
            username,
            user_id,
            bump,
        )?;
        Ok(())
    }

    /// Update the username of an identity
    /// This is only allowed by the account holder
    pub fn update_username(ctx: Context<UpdateUsername>, username: String) -> Result<()> {
        let identity = &mut ctx.accounts.identity;
        identity.update_username(username);
        Ok(())
    }

    pub fn transfer_ownership(ctx: Context<TransferOwnership>) -> Result<()> {
        let identity = &mut ctx.accounts.identity;
        identity.address = ctx.accounts.account_holder_new.key();
        Ok(())
    }

    pub fn delete_identity(_ctx: Context<DeleteIdentity>) -> Result<()> {
        Ok(())
    }
}

#[derive(Accounts)]
pub struct Initialize<'info> {
    #[account(mut)]
    pub protocol_owner: Signer<'info>,

    #[account(
        init,
        seeds = [IDENTITY_SEED.as_bytes()],
        payer=protocol_owner,
        space= 8 + size_of::<IdentityProgram>(),
        bump,
    )]
    pub identity_program: Account<'info, IdentityProgram>,

    pub system_program: Program<'info, System>,
}

#[derive(Accounts)]
#[instruction(social: String, username: String,user_id: u64)]
pub struct CreateIdentity<'info> {
    /// the web3 address owner
    #[account(mut)]
    account_holder: Signer<'info>,

    /// the protocol owner is needed
    /// to verify that the account holder is allowed to create identities
    /// The protocol owner is responsible for the link being valid
    #[account(mut)]
    protocol_owner: Signer<'info>,

    #[account(
        constraint = protocol_owner.key() == identity_program.protocol_owner @ ErrorCodes::ProtocolOwnerNotOwner
    )]
    identity_program: Account<'info, IdentityProgram>,

    #[account(
        init,
        payer=account_holder,
        seeds = [
            // seeds permits one web2 <-> web3 per social account 
            IDENTITY_SEED.as_bytes(),
            social.as_bytes(),
            user_id.to_le_bytes().as_ref()
        ],
        space= 8 + 32 +4 +32 +8+4+32+1+32,
        bump,
    )]
    identity: Account<'info, Identity>,

    system_program: Program<'info, System>,
}

/// Update the username of the identity. This does not effect
/// the id of the user
#[derive(Accounts)]
pub struct UpdateUsername<'info> {
    #[account(mut)]
    account_holder: Signer<'info>,

    #[account(
        mut,
        seeds = [
            IDENTITY_SEED.as_bytes(),
            identity.social_raw.as_bytes(),
            identity.user_id.to_le_bytes().as_ref()
        ],
        bump = identity.bump,
        constraint = account_holder.key() == identity.address @ ErrorCodes::SignerNotOwner
    )]
    identity: Account<'info, Identity>,

    system_program: Program<'info, System>,
}

/// Transfer the ownership of the identity to a new address by
/// co signing with the new key.
#[derive(Accounts)]
pub struct TransferOwnership<'info> {
    #[account(mut)]
    account_holder_curr: Signer<'info>,

    #[account(mut)]
    account_holder_new: Signer<'info>,

    #[account(
        mut,
        seeds = [
            IDENTITY_SEED.as_bytes(),
            identity.social_raw.as_bytes(),
            identity.user_id.to_le_bytes().as_ref()
        ],
        bump = identity.bump,
        constraint = account_holder_curr.key() == identity.address @ ErrorCodes::SignerNotOwner
    )]
    identity: Account<'info, Identity>,

    system_program: Program<'info, System>,
}

/// Delete the identity is the same as closing the account
#[derive(Accounts)]
pub struct DeleteIdentity<'info> {
    #[account(mut)]
    account_holder: Signer<'info>,

    #[account(
        mut,
        close = account_holder,
        seeds = [
            IDENTITY_SEED.as_bytes(),
            identity.social_raw.as_bytes(),
            identity.user_id.to_le_bytes().as_ref()
        ],
        bump = identity.bump,
        constraint = account_holder.key() == identity.address @ ErrorCodes::SignerNotOwner
    )]
    identity: Account<'info, Identity>,

    system_program: Program<'info, System>,
}

#[account]
pub struct IdentityProgram {
    pub protocol_owner: Pubkey,
    pub bump: u8,
}

/// The identity is the account that is used to link
/// a web2 account to a web3 account
///
/// The layout is created to make it easy to use memcmp
/// to query the accounts
#[account]
pub struct Identity {
    // address of the account holder
    pub address: Pubkey, // 32 bytes
    pub social: String,  // 4+32 bytes

    /// the id of the user on the social media
    /// this is immutable
    pub user_id: u64, // 8 bytes

    /// the username of the user on the social media
    /// this is mutable
    pub username: String, // 4+32 bytes

    /// the bump is used to generate the address
    pub bump: u8, // 1 byte

    /// ending bytes used as seed, still assume that it's less than 32 bytes
    pub social_raw: String, // max 32 bytes
}

impl Identity {
    pub fn get_u32string(&self, str: String) -> String {
        if str.len() > 32 {
            panic!("Username is too long")
        }
        // allocate a 32 bytes array and fill it with the username
        let array_of_zeros = vec![0u8; 32 - str.len()];
        return str + std::str::from_utf8(&array_of_zeros).unwrap();
    }

    pub fn init(
        &mut self,
        address: Pubkey,
        social_raw: String,
        username_raw: String,
        user_id: u64,
        bump: &u8,
    ) -> Result<()> {
        self.social_raw = social_raw.clone();
        let username = self.get_u32string(username_raw);
        let social = self.get_u32string(social_raw);
        self.username = username;
        self.address = address;
        self.social = social;
        self.user_id = user_id;
        self.bump = *bump;

        Ok(())
    }

    pub fn update_username(&mut self, username: String) {
        self.username = self.get_u32string(username);
    }
}
