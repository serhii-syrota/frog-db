import * as apiGen from './apiCodegen';
import { FutureResE } from './types';
import { to } from 'await-to-js';

class Api {
  constructor(private readonly _url: string) {}
  private readonly _client = apiGen.DefaultApiFp(
    new apiGen.Configuration({ basePath: this._url })
  );

  public async dbSchema(): FutureResE<DbSchema> {
    const [err, res] = await to((await this._client.dbSchema())());
    if (err) {
      return [err, null];
    }
    return [null, res.data as DbSchema];
  }

  public async createTable(
    tableSchema: Required<apiGen.TableSchema>
  ): FutureResE<apiGen.Info> {
    const [err, res] = await to(
      (
        await this._client.createTable(tableSchema)
      )()
    );
    if (err) {
      return [err, null];
    }
    return [null, res.data];
  }

  public async dropTable(name: string): FutureResE<apiGen.Info> {
    const [err, res] = await to(
      (
        await this._client.deleteTable(name)
      )()
    );
    if (err) {
      return [err, null];
    }
    return [null, res.data];
  }

  public async deleteDuplicateRows(
    name: string
  ): FutureResE<apiGen.Info> {
    const [err, res] = await to(
      (
        await this._client.deleteDuplicateRows(name)
      )()
    );
    if (err) {
      return [err, null];
    }
    return [null, res.data];
  }

  public async deleteRows(
    name: string,
    conditions: Record<string, unknown>
  ): FutureResE<apiGen.Info> {
    const [err, res] = await to(
      (
        await this._client.deleteRows(name, conditions)
      )()
    );
    if (err) {
      return [err, null];
    }
    return [null, res.data];
  }

  public async insertRows(
    name: string,
    rows: Record<string, unknown>[]
  ): FutureResE<apiGen.Info> {
    const [err, res] = await to(
      (
        await this._client.insertRows(name, rows)
      )()
    );
    if (err) {
      return [err, null];
    }
    return [null, res.data];
  }

  public async selectRows(
    name: string,
    conditions: apiGen.SelectBody
  ): FutureResE<Record<string, unknown>[]> {
    const [err, res] = await to(
      (
        await this._client.selectRows(name, conditions)
      )()
    );
    if (err) {
      return [err, null];
    }
    return [null, res.data];
  }

  public async updateRows(
    name: string,
    data: apiGen.UpdateBody
  ): FutureResE<apiGen.Info> {
    const [err, res] = await to(
      (
        await this._client.updateRows(name, data)
      )()
    );
    if (err) {
      return [err, null];
    }
    return [null, res.data];
  }
}

export const api = new Api('http://localhost:8080');

export type DbSchema = Required<apiGen.TableSchema>[];
