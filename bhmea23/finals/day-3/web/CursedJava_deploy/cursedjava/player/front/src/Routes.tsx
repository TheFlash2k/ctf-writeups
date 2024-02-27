import {
    Routes,
    Route,
} from 'react-router-dom'
import { Dashboard } from './Dashboard'
import Home from './Home'
import { Login } from "./Login";
import { Register } from "./Register";


export default function Routing() {
    return (
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/register" element={<Register />} />
        <Route path="/login" element={<Login />} />
      </Routes>
    );
}
