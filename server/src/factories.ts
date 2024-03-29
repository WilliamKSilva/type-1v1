import { PoolClient } from "pg";
import { GameController } from "./controllers/game_controller";
import { GameRepository } from "./repositories/game_repository";
import { GameService } from "./services/game_service";
import { Server } from "socket.io";
import { GameWebSocketService } from "./services/game_websocket_service";
import { UserRepository } from "./repositories/user_repository";
import { UserService } from "./services/user_service";
import { UserController } from "./controllers/user_controller";

export const NewGameFactory = (pool: PoolClient): GameController => {
  const gameRepository = new GameRepository(pool);
  const gameService = new GameService(gameRepository);
  const gameController = new GameController(gameService);

  return gameController;
};

export const NewGameWebSocketFactory = (io: Server): GameWebSocketService => {
  const gameWebSocketService = new GameWebSocketService(io);

  return gameWebSocketService;
};

export const NewUserFactory = (pool: PoolClient): UserController => {
  const userRepository = new UserRepository(pool);
  const userService = new UserService(userRepository);

  const userController = new UserController(userService);

  return userController;
};
