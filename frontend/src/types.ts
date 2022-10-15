export type FutureResE<T = never, U = Error> = [T] extends [never]
  ? Promise<U | null>
  : Promise<[U, null] | [null, T]>;
