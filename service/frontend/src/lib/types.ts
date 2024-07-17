export type Devenv = {
  id: string;
  public: boolean;
  name: string;
  buildCmd: string;
  runCmd: string;
  created: string;
  updated: string;
};

export type GetUserResponse = {
  id: string;
  username: string;
  created: string;
  updated: string;
  deleted: string;
};

export type CreateDevenvRequest = {
  name: string;
  buildCmd: string;
  runCmd: string;
};

export type PatchDevenvRequest = {
  name?: string;
  buildCmd?: string;
  runCmd?: string;
};

export type CreateDevenvResponse = {
  devenvUuid: string;
};
