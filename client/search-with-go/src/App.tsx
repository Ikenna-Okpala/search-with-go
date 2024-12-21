import Go from "./assets/go-mascot.svg";
import { IoSearch } from "react-icons/io5";
import { Input } from "./components/ui/input";
import { useNavigate } from "react-router-dom";
import { useState } from "react";

function App() {
  const navigate = useNavigate();

  const [search, setSearch] = useState("");

  const onEnterPressed = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      navigate("/search", {
        state: {
          search,
        },
      });
    }
  };

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearch(e.target.value);
  };
  return (
    <div className="w-screen h-screen flex flex-col justify-center items-center gap-y-8">
      <div className="flex flex-row items-center">
        <span className="text-8xl font-medium text-white">Search With</span>
        <img src={Go} alt="Go mascot" width={150} height={150} />
      </div>

      <div className="flex flex-row gap-x-3 w-2/4 items-center bg-gray200 rounded-full px-5 shadow-lg">
        <IoSearch className="text-4xl text-gray100" />
        <Input
          type="text"
          className="h-[75px] bg-transparent border-none text-white text-2xl font-medium"
          onKeyDown={onEnterPressed}
          onChange={onChange}
          value={search}
        />
      </div>
    </div>
  );
}

export default App;
