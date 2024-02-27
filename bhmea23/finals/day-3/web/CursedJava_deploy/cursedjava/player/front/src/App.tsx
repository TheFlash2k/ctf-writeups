
import { HeaderMiddle } from "./Navbar";
import Routing from "./Routes";

function App() {
  return (
    <div>
      <HeaderMiddle
        links={[
          {
            link: "/",
            label: "Home",
          },
          {
            link: "/dashboard",
            label: "Dashboard",
          },
          {
            link: "/dashboard",
            label: "Profile",
          },
          {
            link: "/register",
            label: "Register",
          },
          {
            link: "/login",
            label: "Login",
          },
        ]}
      />
      <Routing />

    </div>
  );
}

export default App;
