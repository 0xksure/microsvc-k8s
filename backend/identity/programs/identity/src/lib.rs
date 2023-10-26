use anchor_lang::prelude::*;
use std::mem::size_of;

declare_id!("3rQketG7pSopHE1APQKZu1BQofanqbCBP7spZ4CBGrUm");

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
        user_id: u32,
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
#[instruction(social: String, username: String,user_id: u32)]
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
        space= 8 + size_of::<Identity>(),
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

#[derive(AnchorSerialize, AnchorDeserialize, Clone, Debug)]
pub enum Social {
    Facebook,
    Twitter,
    Instagram,
    LinkedIn,
    Github,
    Website,
    Email,
}
impl Social {
    pub fn to_u8(&self) -> u8 {
        match self {
            Social::Facebook => 0,
            Social::Twitter => 1,
            Social::Instagram => 2,
            Social::LinkedIn => 3,
            Social::Github => 4,
            Social::Website => 5,
            Social::Email => 6,
        }
    }

    pub fn from_string(social: &str) -> Social {
        match social {
            "facebook" => Social::Facebook,
            "twitter" => Social::Twitter,
            "instagram" => Social::Instagram,
            "linkedin" => Social::LinkedIn,
            "github" => Social::Github,
            "website" => Social::Website,
            "email" => Social::Email,
            _ => panic!("Invalid social media"),
        }
    }
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
    pub social: Social,  // 1 byte

    /// the id of the user on the social media
    /// this is immutable
    pub user_id: u32, // 4 bytes

    /// the username of the user on the social media
    /// this is mutable
    pub username: Vec<u8>, // 32 bytes

    /// the bump is used to generate the address
    pub bump: u8,

    pub social_raw: String,
}

impl Identity {
    pub fn init(
        &mut self,
        address: Pubkey,
        social_raw: String,
        username_raw: String,
        user_id: u32,
        bump: &u8,
    ) -> Result<()> {
        let username = username_raw.as_bytes().to_vec();
        if username.len() > 32 {
            return Err(ErrorCodes::UsernameTooLong.into());
        }
        self.username = username;
        self.address = address;
        self.social = Social::from_string(&social_raw);
        self.social_raw = social_raw;
        self.user_id = user_id;
        self.bump = *bump;

        Ok(())
    }

    pub fn update_username(&mut self, username: String) {
        if username.len() > 32 {
            panic!("Social name is too long")
        }
        self.username = username.as_bytes().to_vec();
    }
}
