export type Identity = {
    "version": "0.1.0";
    "name": "identity";
    "instructions": [
        {
            "name": "initialize";
            "docs": [
                "Initialize the identity program",
                "Sets the protocol owner that is allowed to create identities"
            ];
            "accounts": [
                {
                    "name": "protocolOwner";
                    "isMut": true;
                    "isSigner": true;
                },
                {
                    "name": "identityProgram";
                    "isMut": true;
                    "isSigner": false;
                },
                {
                    "name": "systemProgram";
                    "isMut": false;
                    "isSigner": false;
                }
            ];
            "args": [];
        },
        {
            "name": "createIdentity";
            "docs": [
                "Create a new web2 <-> web3 identity"
            ];
            "accounts": [
                {
                    "name": "accountHolder";
                    "isMut": true;
                    "isSigner": true;
                    "docs": [
                        "the web3 address owner"
                    ];
                },
                {
                    "name": "protocolOwner";
                    "isMut": true;
                    "isSigner": true;
                    "docs": [
                        "the protocol owner is needed",
                        "to verify that the account holder is allowed to create identities",
                        "The protocol owner is responsible for the link being valid"
                    ];
                },
                {
                    "name": "identityProgram";
                    "isMut": false;
                    "isSigner": false;
                },
                {
                    "name": "identity";
                    "isMut": true;
                    "isSigner": false;
                },
                {
                    "name": "systemProgram";
                    "isMut": false;
                    "isSigner": false;
                }
            ];
            "args": [
                {
                    "name": "social";
                    "type": "string";
                },
                {
                    "name": "username";
                    "type": "string";
                },
                {
                    "name": "userId";
                    "type": "u32";
                }
            ];
        },
        {
            "name": "updateUsername";
            "docs": [
                "Update the username of an identity",
                "This is only allowed by the account holder"
            ];
            "accounts": [
                {
                    "name": "accountHolder";
                    "isMut": true;
                    "isSigner": true;
                },
                {
                    "name": "identity";
                    "isMut": true;
                    "isSigner": false;
                },
                {
                    "name": "systemProgram";
                    "isMut": false;
                    "isSigner": false;
                }
            ];
            "args": [
                {
                    "name": "username";
                    "type": "string";
                }
            ];
        },
        {
            "name": "transferOwnership";
            "accounts": [
                {
                    "name": "accountHolderCurr";
                    "isMut": true;
                    "isSigner": true;
                },
                {
                    "name": "accountHolderNew";
                    "isMut": true;
                    "isSigner": true;
                },
                {
                    "name": "identity";
                    "isMut": true;
                    "isSigner": false;
                },
                {
                    "name": "systemProgram";
                    "isMut": false;
                    "isSigner": false;
                }
            ];
            "args": [];
        },
        {
            "name": "deleteIdentity";
            "accounts": [
                {
                    "name": "accountHolder";
                    "isMut": true;
                    "isSigner": true;
                },
                {
                    "name": "identity";
                    "isMut": true;
                    "isSigner": false;
                },
                {
                    "name": "systemProgram";
                    "isMut": false;
                    "isSigner": false;
                }
            ];
            "args": [];
        }
    ];
    "accounts": [
        {
            "name": "identityProgram";
            "type": {
                "kind": "struct";
                "fields": [
                    {
                        "name": "protocolOwner";
                        "type": "publicKey";
                    },
                    {
                        "name": "bump";
                        "type": "u8";
                    }
                ];
            };
        },
        {
            "name": "identity";
            "docs": [
                "The identity is the account that is used to link",
                "a web2 account to a web3 account",
                "",
                "The layout is created to make it easy to use memcmp",
                "to query the accounts"
            ];
            "type": {
                "kind": "struct";
                "fields": [
                    {
                        "name": "address";
                        "type": "publicKey";
                    },
                    {
                        "name": "social";
                        "type": {
                            "defined": "Social";
                        };
                    },
                    {
                        "name": "userId";
                        "docs": [
                            "the id of the user on the social media",
                            "this is immutable"
                        ];
                        "type": "u32";
                    },
                    {
                        "name": "username";
                        "docs": [
                            "the username of the user on the social media",
                            "this is mutable"
                        ];
                        "type": "bytes";
                    },
                    {
                        "name": "bump";
                        "docs": [
                            "the bump is used to generate the address"
                        ];
                        "type": "u8";
                    },
                    {
                        "name": "socialRaw";
                        "type": "string";
                    }
                ];
            };
        }
    ];
    "types": [
        {
            "name": "Social";
            "type": {
                "kind": "enum";
                "variants": [
                    {
                        "name": "Facebook";
                    },
                    {
                        "name": "Twitter";
                    },
                    {
                        "name": "Instagram";
                    },
                    {
                        "name": "LinkedIn";
                    },
                    {
                        "name": "Github";
                    },
                    {
                        "name": "Website";
                    },
                    {
                        "name": "Email";
                    }
                ];
            };
        }
    ];
    "errors": [
        {
            "code": 6000;
            "name": "ProtocolOwnerNotOwner";
            "msg": "The protocol owner is not the owner of the account";
        },
        {
            "code": 6001;
            "name": "SignerNotOwner";
            "msg": "The signer is not the owner of the account";
        },
        {
            "code": 6002;
            "name": "UsernameTooLong";
            "msg": "The username is too long. Should be max 32bytes";
        }
    ];
};
export declare const IDL: Identity;
//# sourceMappingURL=identity.d.ts.map