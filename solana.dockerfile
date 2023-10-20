FROM ghcr.io/wormhole-foundation/solana:1.10.31@sha256:d31e8db926a1d3fbaa9d9211d9979023692614b7b64912651aba0383e8c01bad AS solana


# Support additional root CAs
COPY cert.pem* /certs/
# Debian
RUN if [ -e /certs/cert.pem ]; then cp /certs/cert.pem /etc/ssl/certs/ca-certificates.crt; fi

# Add bridge contract sources
WORKDIR /usr/src/bridge

COPY programs programs
COPY Anchor.toml Anchor.toml
COPY Cargo.toml Cargo.toml
COPY Cargo.lock Cargo.lock
COPY package.json package.json
COPY package-lock.json package-lock.json
COPY yarn.lock yarn.lock
COPY target/deploy/identity.so target/deploy/identity.so

ENV RUST_LOG="solana_runtime::system_instruction_processor=trace,solana_runtime::message_processor=trace,solana_bpf_loader=debug,solana_rbpf=debug"
ENV RUST_BACKTRACE=1


FROM solana AS builder

RUN mkdir -p /opt/solana/deps

ARG EMITTER_ADDRESS="11111111111111111111111111111115"

# RUN mkdir -p /opt/solana/deps
# RUN rustup toolchain add 1.73
# # Build Wormhole Solana programs

# RUN rm -rf ~/.cargo/registry && cargo build
# RUN rm -rf /usr/local/cargo/registry && cargo build-bpf --manifest-path "programs/identity/Cargo.toml" 
# RUN solana --version
# RUN anchor --version
# RUN rustc --version
# RUN solana build-bpf  
RUN cp target/deploy/identity.so /opt/solana/deps/identity.so 

# External imports 
#cp external/mpl_token_metadata.so /opt/solana/deps/mpl_token_metadata.so

FROM scratch AS export-stage
COPY --from=builder /opt/solana/deps /