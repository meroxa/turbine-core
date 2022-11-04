// package: turbine_core
// file: turbine.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";

export class InitRequest extends jspb.Message {
  getAppname(): string;
  setAppname(value: string): void;

  getConfigfilepath(): string;
  setConfigfilepath(value: string): void;

  getLanguage(): InitRequest.LanguageMap[keyof InitRequest.LanguageMap];
  setLanguage(value: InitRequest.LanguageMap[keyof InitRequest.LanguageMap]): void;

  getGitsha(): string;
  setGitsha(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InitRequest.AsObject;
  static toObject(includeInstance: boolean, msg: InitRequest): InitRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: InitRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InitRequest;
  static deserializeBinaryFromReader(message: InitRequest, reader: jspb.BinaryReader): InitRequest;
}

export namespace InitRequest {
  export type AsObject = {
    appname: string,
    configfilepath: string,
    language: InitRequest.LanguageMap[keyof InitRequest.LanguageMap],
    gitsha: string,
  }

  export interface LanguageMap {
    GOLANG: 0;
    PYTHON: 1;
    JAVASCRIPT: 2;
    RUBY: 3;
  }

  export const Language: LanguageMap;
}

export class NameOrUUID extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getUuid(): string;
  setUuid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NameOrUUID.AsObject;
  static toObject(includeInstance: boolean, msg: NameOrUUID): NameOrUUID.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: NameOrUUID, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NameOrUUID;
  static deserializeBinaryFromReader(message: NameOrUUID, reader: jspb.BinaryReader): NameOrUUID;
}

export namespace NameOrUUID {
  export type AsObject = {
    name: string,
    uuid: string,
  }
}

export class Resource extends jspb.Message {
  getUuid(): string;
  setUuid(value: string): void;

  getName(): string;
  setName(value: string): void;

  getType(): string;
  setType(value: string): void;

  getDirection(): Resource.DirectionMap[keyof Resource.DirectionMap];
  setDirection(value: Resource.DirectionMap[keyof Resource.DirectionMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Resource.AsObject;
  static toObject(includeInstance: boolean, msg: Resource): Resource.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Resource, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Resource;
  static deserializeBinaryFromReader(message: Resource, reader: jspb.BinaryReader): Resource;
}

export namespace Resource {
  export type AsObject = {
    uuid: string,
    name: string,
    type: string,
    direction: Resource.DirectionMap[keyof Resource.DirectionMap],
  }

  export interface DirectionMap {
    DIRECTION_SOURCE: 0;
    DIRECTION_DESTINATION: 1;
  }

  export const Direction: DirectionMap;
}

export class Collection extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getStream(): string;
  setStream(value: string): void;

  clearRecordsList(): void;
  getRecordsList(): Array<Record>;
  setRecordsList(value: Array<Record>): void;
  addRecords(value?: Record, index?: number): Record;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Collection.AsObject;
  static toObject(includeInstance: boolean, msg: Collection): Collection.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Collection, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Collection;
  static deserializeBinaryFromReader(message: Collection, reader: jspb.BinaryReader): Collection;
}

export namespace Collection {
  export type AsObject = {
    name: string,
    stream: string,
    recordsList: Array<Record.AsObject>,
  }
}

export class Record extends jspb.Message {
  getKey(): string;
  setKey(value: string): void;

  getValue(): Uint8Array | string;
  getValue_asU8(): Uint8Array;
  getValue_asB64(): string;
  setValue(value: Uint8Array | string): void;

  hasTimestamp(): boolean;
  clearTimestamp(): void;
  getTimestamp(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setTimestamp(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Record.AsObject;
  static toObject(includeInstance: boolean, msg: Record): Record.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Record, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Record;
  static deserializeBinaryFromReader(message: Record, reader: jspb.BinaryReader): Record;
}

export namespace Record {
  export type AsObject = {
    key: string,
    value: Uint8Array | string,
    timestamp?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class ReadCollectionRequest extends jspb.Message {
  hasResource(): boolean;
  clearResource(): void;
  getResource(): Resource | undefined;
  setResource(value?: Resource): void;

  getCollection(): string;
  setCollection(value: string): void;

  hasConfigs(): boolean;
  clearConfigs(): void;
  getConfigs(): ResourceConfigs | undefined;
  setConfigs(value?: ResourceConfigs): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReadCollectionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ReadCollectionRequest): ReadCollectionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ReadCollectionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReadCollectionRequest;
  static deserializeBinaryFromReader(message: ReadCollectionRequest, reader: jspb.BinaryReader): ReadCollectionRequest;
}

export namespace ReadCollectionRequest {
  export type AsObject = {
    resource?: Resource.AsObject,
    collection: string,
    configs?: ResourceConfigs.AsObject,
  }
}

export class WriteCollectionRequest extends jspb.Message {
  hasResource(): boolean;
  clearResource(): void;
  getResource(): Resource | undefined;
  setResource(value?: Resource): void;

  hasCollection(): boolean;
  clearCollection(): void;
  getCollection(): Collection | undefined;
  setCollection(value?: Collection): void;

  getTargetcollection(): string;
  setTargetcollection(value: string): void;

  hasConfigs(): boolean;
  clearConfigs(): void;
  getConfigs(): ResourceConfigs | undefined;
  setConfigs(value?: ResourceConfigs): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WriteCollectionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: WriteCollectionRequest): WriteCollectionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: WriteCollectionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WriteCollectionRequest;
  static deserializeBinaryFromReader(message: WriteCollectionRequest, reader: jspb.BinaryReader): WriteCollectionRequest;
}

export namespace WriteCollectionRequest {
  export type AsObject = {
    resource?: Resource.AsObject,
    collection?: Collection.AsObject,
    targetcollection: string,
    configs?: ResourceConfigs.AsObject,
  }
}

export class ResourceConfigs extends jspb.Message {
  clearResourceconfigList(): void;
  getResourceconfigList(): Array<ResourceConfig>;
  setResourceconfigList(value: Array<ResourceConfig>): void;
  addResourceconfig(value?: ResourceConfig, index?: number): ResourceConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResourceConfigs.AsObject;
  static toObject(includeInstance: boolean, msg: ResourceConfigs): ResourceConfigs.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ResourceConfigs, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResourceConfigs;
  static deserializeBinaryFromReader(message: ResourceConfigs, reader: jspb.BinaryReader): ResourceConfigs;
}

export namespace ResourceConfigs {
  export type AsObject = {
    resourceconfigList: Array<ResourceConfig.AsObject>,
  }
}

export class ResourceConfig extends jspb.Message {
  getField(): string;
  setField(value: string): void;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResourceConfig.AsObject;
  static toObject(includeInstance: boolean, msg: ResourceConfig): ResourceConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ResourceConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResourceConfig;
  static deserializeBinaryFromReader(message: ResourceConfig, reader: jspb.BinaryReader): ResourceConfig;
}

export namespace ResourceConfig {
  export type AsObject = {
    field: string,
    value: string,
  }
}

export class Process extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getType(): Process.TypeMap[keyof Process.TypeMap];
  setType(value: Process.TypeMap[keyof Process.TypeMap]): void;

  getImage(): string;
  setImage(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Process.AsObject;
  static toObject(includeInstance: boolean, msg: Process): Process.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Process, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Process;
  static deserializeBinaryFromReader(message: Process, reader: jspb.BinaryReader): Process;
}

export namespace Process {
  export type AsObject = {
    name: string,
    type: Process.TypeMap[keyof Process.TypeMap],
    image: string,
  }

  export interface TypeMap {
    GO: 0;
    NODE: 1;
    PYTHON: 2;
    DOCKER: 3;
  }

  export const Type: TypeMap;
}

export class ProcessCollectionRequest extends jspb.Message {
  hasProcess(): boolean;
  clearProcess(): void;
  getProcess(): Process | undefined;
  setProcess(value?: Process): void;

  hasCollection(): boolean;
  clearCollection(): void;
  getCollection(): Collection | undefined;
  setCollection(value?: Collection): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProcessCollectionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ProcessCollectionRequest): ProcessCollectionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProcessCollectionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProcessCollectionRequest;
  static deserializeBinaryFromReader(message: ProcessCollectionRequest, reader: jspb.BinaryReader): ProcessCollectionRequest;
}

export namespace ProcessCollectionRequest {
  export type AsObject = {
    process?: Process.AsObject,
    collection?: Collection.AsObject,
  }
}

export class Secret extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Secret.AsObject;
  static toObject(includeInstance: boolean, msg: Secret): Secret.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Secret, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Secret;
  static deserializeBinaryFromReader(message: Secret, reader: jspb.BinaryReader): Secret;
}

export namespace Secret {
  export type AsObject = {
    name: string,
    value: string,
  }
}

export class ListFunctionsResponse extends jspb.Message {
  clearFunctionsList(): void;
  getFunctionsList(): Array<string>;
  setFunctionsList(value: Array<string>): void;
  addFunctions(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFunctionsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListFunctionsResponse): ListFunctionsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListFunctionsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFunctionsResponse;
  static deserializeBinaryFromReader(message: ListFunctionsResponse, reader: jspb.BinaryReader): ListFunctionsResponse;
}

export namespace ListFunctionsResponse {
  export type AsObject = {
    functionsList: Array<string>,
  }
}

export class ResourceWithCollection extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getCollection(): string;
  setCollection(value: string): void;

  getDirection(): ResourceWithCollection.DirectionMap[keyof ResourceWithCollection.DirectionMap];
  setDirection(value: ResourceWithCollection.DirectionMap[keyof ResourceWithCollection.DirectionMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResourceWithCollection.AsObject;
  static toObject(includeInstance: boolean, msg: ResourceWithCollection): ResourceWithCollection.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ResourceWithCollection, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResourceWithCollection;
  static deserializeBinaryFromReader(message: ResourceWithCollection, reader: jspb.BinaryReader): ResourceWithCollection;
}

export namespace ResourceWithCollection {
  export type AsObject = {
    name: string,
    collection: string,
    direction: ResourceWithCollection.DirectionMap[keyof ResourceWithCollection.DirectionMap],
  }

  export interface DirectionMap {
    DIRECTION_SOURCE: 0;
    DIRECTION_DESTINATION: 1;
  }

  export const Direction: DirectionMap;
}

export class ListResourcesResponse extends jspb.Message {
  clearResourcesList(): void;
  getResourcesList(): Array<ResourceWithCollection>;
  setResourcesList(value: Array<ResourceWithCollection>): void;
  addResources(value?: ResourceWithCollection, index?: number): ResourceWithCollection;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListResourcesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListResourcesResponse): ListResourcesResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListResourcesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListResourcesResponse;
  static deserializeBinaryFromReader(message: ListResourcesResponse, reader: jspb.BinaryReader): ListResourcesResponse;
}

export namespace ListResourcesResponse {
  export type AsObject = {
    resourcesList: Array<ResourceWithCollection.AsObject>,
  }
}

