use anchor_lang::prelude::*;

declare_id!("3rQketG7pSopHE1APQKZu1BQofanqbCBP7spZ4CBGrUm");

#[program]
pub mod identity {
    use super::*;

    pub fn initialize(ctx: Context<Initialize>) -> Result<()> {
        Ok(())
    }
}

#[derive(Accounts)]
pub struct Initialize {}
