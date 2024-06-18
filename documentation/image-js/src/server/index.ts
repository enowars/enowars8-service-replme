import createExpressApp from "./router";

export default function init() {
  const app = createExpressApp();

  const port = 3000;
  const host = '0.0.0.0';

  console.log('App listening to http://' + host + ':' + port);
  app.listen(port, host, 0);
}

