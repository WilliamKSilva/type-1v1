export enum GameType {
  regular = "regular",
  fast = "fast",
}

export class Game {
  uuid: string = "";
  name: string = "";
  type: GameType = GameType.regular;
  created_at: Date = new Date();
  updated_at: Date = new Date();

  constructor(
    uuid: string,
    name: string,
    type: GameType,
    created_at: Date,
    updated_at: Date
  ) {
    this.uuid = uuid;
    this.name = name;
    this.type = type;
    this.created_at = created_at;
    this.updated_at = updated_at;
  }
}

// TODO, update this to match the two types of game, fast and regular
export type PlayerStateData = {
  name: string;
  textState: string;
};
