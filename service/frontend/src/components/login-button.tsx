import { Button } from "./ui/button";
import { navigate } from "@/actions/navigate";

export function LoginButton() {
  return (
    <Button variant="outline" onClick={() => navigate("/login")}>
      LOGIN
    </Button>
  );
}
