import axios from 'axios'
import { useEffect, useMemo } from "react";
import { useMutation, useQuery } from "react-query";
import { useNavigate } from "react-router-dom";

const url = "/api";

const client = axios.create({
  baseURL: url,
});

function simpleErrorParser(e: string) {
  if (e.includes("java.lang")) {
    return e.split(":")[1];
  }
  return e;
}

const register = (username: string, password: string) =>
  client.post("/user/register", { username, password }).then((res) => res.data);

const login = (username: string, password: string) =>
  client
    .post(
      "/user/login",
      { username, password },
      {
        withCredentials: true,
      }
    )
    .then((res) => res.data);

const redeem = (code: string) =>
  client.get("/coupon/use/" + code).then((res) => res.data);

const getMyData = () => client.get("/user").then((res) => res.data);

const getFlag = () => client.get("/user/flag").then((res) => res.data);

const logout = () => client.post("/user/logout").then((res) => res.data);

export const useGetMyData = () => {
  const myData = useQuery("user", getMyData);
  return myData;
};

export const useIsLoggedIn = () => {
  const myData = useQuery("user", getMyData);
  return !!myData.data;
};

export const useGetFlag = () => {
  return useQuery("flag", getFlag);
};

export const useLogout = () => {
  const myData = useQuery("user", getMyData);
  const push = useNavigate();
  const { mutate } = useMutation(() => logout(), {
    onSuccess: () => {
      myData.refetch();
      push("/");
    },
  });

  return mutate;
};

export const useRegister = () => {
  const myData = useQuery("user", getMyData);
  const push = useNavigate();
  const { mutate, error: e } = useMutation(
    (data: { username: string; password: string }) =>
      register(data.username, data.password),
    {
      onSuccess: () => {
        myData.refetch();
        push("/login");
      },
    }
  );

  const error = useMemo(() => {
    if (!e) {
      return null;
    }
    const ee = e as any;
    if (ee.response) {
      return simpleErrorParser(ee.response.data);
    }

    return ee.message;
  }, [e]);

  return [mutate, error] as const;
};

export const useLogin = () => {
  const myData = useQuery("user", getMyData);
  const push = useNavigate();
  const { mutate, error: e } = useMutation(
    (data: { username: string; password: string }) =>
      login(data.username, data.password),
    {
      onSuccess: () => {
        myData.refetch();
        push("/");
      },
    }
  );

  const error = useMemo(() => {
    if (!e) {
      return null;
    }
    const ee = e as any;
    if (ee.response) {
      return simpleErrorParser(ee.response.data);
    }

    return ee.message;
  }, [e]);

  return [mutate, error] as const;
};

export const useRedeem = () => {
  const { mutate, error: e } = useMutation((code: string) => redeem(code), {
    onSuccess: () => {
      alert("Coupon redeemed!");
    },
  });

  const error = useMemo(() => {
    if (!e) {
      return null;
    }
    const ee = e as any;
    if (ee.response) {
      return simpleErrorParser(ee.response.data);
    }

    return ee.message;
  }, [e]);

  return [mutate, error] as const;
};