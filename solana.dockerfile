FROM rust:1.73-buster as rust

RUN sh -c "$(curl -sSfL https://release.solana.com/v1.17.10/install)"

ENV PATH="/root/.local/share/solana/install/active_release/bin:$PATH"




# The strip shell script downloads criterion the first time it runs so cache it here as well.
# RUN touch /tmp/foo.so && \
#     /root/.local/share/solana/install/active_release/bin/sdk/bpf/scripts/strip.sh /tmp/foo.so /tmp/bar.so || \
#     rm /tmp/foo.so

FROM rust AS solana-setup

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

ENV RUST_LOG="solana_runtime::system_instruction_processor=trace,solana_runtime::message_processor=trace,solana_bpf_loader=debug,solana_rbpf=debug"
ENV RUST_BACKTRACE=1
RUN cargo install --git https://github.com/project-serum/anchor --tag v0.29.0 anchor-cli 


FROM solana-setup AS builder

RUN mkdir -p /opt/solana/deps
RUN rustup toolchain add 1.73
# Build Wormhole Solana programs
RUN --mount=type=cache,target=target,id=build 
RUN --mount=type=cache,target=/usr/local/cargo/registry,id=cargo_registry 
RUN sh -c "$(curl -sSfL https://release.solana.com/v1.17.1/install)"
ENV PATH="/root/.local/share/solana/install/active_release/bin:$PATH"

RUN solana --version
RUN anchor --version
RUN rustc --version
RUN anchor build  
RUN cp target/deploy/identity.so /opt/solana/deps/identity.so 

# External imports 
#cp external/mpl_token_metadata.so /opt/solana/deps/mpl_token_metadata.so

FROM scratch AS export-stage
COPY --from=builder /opt/solana/deps /