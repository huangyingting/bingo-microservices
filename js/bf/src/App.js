import BsNavbar from "./components/BsNavbar";
import { BrowserRouter, Route, Routes, Navigate } from "react-router-dom";
import Home from './components/Home'
import Stats from './components/Stats'
import ShortUrl from './components/ShortUrl'
import About from './components/About'
import Footer from './components/Footer'
import Unauthenticated from "./components/Unauthenticated";
import PageNotFound from "./components/PageNotFound";

import { AuthenticatedTemplate, UnauthenticatedTemplate } from "@azure/msal-react";

function App() {
  return (
    <div>
      <AuthenticatedTemplate>
        <BrowserRouter>
          <BsNavbar />
          <Routes>
            <Route exact path="/pages/stats" element={<Stats />} />
            <Route exact path="/pages/shorturl" element={<ShortUrl />} />
            <Route exact path="/pages/about" element={<About />} />
            <Route exact path="/" element={<Home />} />
            <Route exact path="/pages" element={<Navigate replace to="/" />} />
            <Route exact path='*' element={<PageNotFound />} />
          </Routes>
          <Footer />
        </BrowserRouter>
      </AuthenticatedTemplate>
      <UnauthenticatedTemplate>
        <Unauthenticated />
      </UnauthenticatedTemplate>
    </div>
  );
}

export default App;
