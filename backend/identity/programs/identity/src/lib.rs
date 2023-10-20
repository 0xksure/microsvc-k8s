use anchor_lang::prelude::*;

declare_id!("3rQketG7pSopHE1APQKZu1BQofanqbCBP7spZ4CBGrUm");

#[program]
pub mod identity {
    use super::*;

    pub fn initialize(ctx: Context<Initialize>) -> Result<()> {
        Ok(())
    }

    pub fn create_identity(ctx: Context<CreateIdentity> ) -> Result<()>{
        Ok(())
    }
}

#[derive(Accounts)]
pub struct Initialize {}

#[derive(Accounts)]
pub struct CreateIdentity<'info> {
    [account(mut)]
    account_holder: Signer<'info>,

    #[account(
        init,
        payer=account_holder,
        seeds = ["identity",account_holder.key().as_ref()],
        space= 8 + 
        bump,
    )]
    identity: Account<'info, Identity>,

    system_program: Program<'info, System>,
}


#[account]
pub struct Identity {
    pub authority: Pubkey,
    pub authority_seeds: Vec<u8>,
    pub name: String,
    pub email: String,
    pub phone: String,
    pub address: String,
    pub created_at: i64,
    pub updated_at: i64,
}