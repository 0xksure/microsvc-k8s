// @generated by protoc-gen-es v1.3.3 with parameter "target=ts"
// @generated from file index.proto (package bounty, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3, protoInt64 } from "@bufbuild/protobuf";

/**
 * @generated from enum bounty.BountySignStatus
 */
export enum BountySignStatus {
  /**
   * @generated from enum value: CREATED = 0;
   */
  CREATED = 0,

  /**
   * @generated from enum value: SIGNED = 1;
   */
  SIGNED = 1,

  /**
   * @generated from enum value: COMPLETED = 2;
   */
  COMPLETED = 2,

  /**
   * @generated from enum value: FAILED_TO_SIGN = 3;
   */
  FAILED_TO_SIGN = 3,

  /**
   * @generated from enum value: CANCELLED = 4;
   */
  CANCELLED = 4,
}
// Retrieve enum metadata with: proto3.getEnumType(BountySignStatus)
proto3.util.setEnumType(BountySignStatus, "bounty.BountySignStatus", [
  { no: 0, name: "CREATED" },
  { no: 1, name: "SIGNED" },
  { no: 2, name: "COMPLETED" },
  { no: 3, name: "FAILED_TO_SIGN" },
  { no: 4, name: "CANCELLED" },
]);

/**
 * @generated from message bounty.BountyMessage
 */
export class BountyMessage extends Message<BountyMessage> {
  /**
   * @generated from field: bounty.BountySignStatus BountySignStatus = 1;
   */
  BountySignStatus = BountySignStatus.CREATED;

  /**
   * @generated from field: int64 Bountyid = 2;
   */
  Bountyid = protoInt64.zero;

  /**
   * @generated from field: string BountyUIAmount = 3;
   */
  BountyUIAmount = "";

  /**
   * @generated from field: string TokenAddress = 4;
   */
  TokenAddress = "";

  /**
   * @generated from field: string CreatorAddress = 5;
   */
  CreatorAddress = "";

  /**
   * @generated from field: int64 InstallationId = 6;
   */
  InstallationId = protoInt64.zero;

  constructor(data?: PartialMessage<BountyMessage>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "bounty.BountyMessage";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "BountySignStatus", kind: "enum", T: proto3.getEnumType(BountySignStatus) },
    { no: 2, name: "Bountyid", kind: "scalar", T: 3 /* ScalarType.INT64 */ },
    { no: 3, name: "BountyUIAmount", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "TokenAddress", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 5, name: "CreatorAddress", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 6, name: "InstallationId", kind: "scalar", T: 3 /* ScalarType.INT64 */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): BountyMessage {
    return new BountyMessage().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): BountyMessage {
    return new BountyMessage().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): BountyMessage {
    return new BountyMessage().fromJsonString(jsonString, options);
  }

  static equals(a: BountyMessage | PlainMessage<BountyMessage> | undefined, b: BountyMessage | PlainMessage<BountyMessage> | undefined): boolean {
    return proto3.util.equals(BountyMessage, a, b);
  }
}

/**
 * @generated from message bounty.LinkerMessage
 */
export class LinkerMessage extends Message<LinkerMessage> {
  /**
   * @generated from field: string Username = 1;
   */
  Username = "";

  /**
   * @generated from field: string UserId = 2;
   */
  UserId = "";

  /**
   * @generated from field: string WalletAddress = 3;
   */
  WalletAddress = "";

  constructor(data?: PartialMessage<LinkerMessage>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "bounty.LinkerMessage";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "Username", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "UserId", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "WalletAddress", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): LinkerMessage {
    return new LinkerMessage().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): LinkerMessage {
    return new LinkerMessage().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): LinkerMessage {
    return new LinkerMessage().fromJsonString(jsonString, options);
  }

  static equals(a: LinkerMessage | PlainMessage<LinkerMessage> | undefined, b: LinkerMessage | PlainMessage<LinkerMessage> | undefined): boolean {
    return proto3.util.equals(LinkerMessage, a, b);
  }
}
