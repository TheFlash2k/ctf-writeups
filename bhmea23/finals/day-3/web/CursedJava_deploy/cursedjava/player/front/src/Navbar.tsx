import { useState } from 'react';
import {
  createStyles,
  Header,
  Group,
  ActionIcon,
  Container,
  Burger,
  rem,
  Button,
  Avatar,
} from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import {
  IconBrandTwitter,
  IconBrandYoutube,
  IconBrandInstagram,
} from "@tabler/icons-react";
import { Link } from "react-router-dom";
import { useIsLoggedIn, useLogout, useGetMyData } from "./api";

const useStyles = createStyles((theme) => ({
  inner: {
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
    height: rem(56),

    [theme.fn.smallerThan("sm")]: {
      justifyContent: "flex-start",
    },
  },

  links: {
    width: rem(500),
    display: "flex",
    alignItems: "center",

    [theme.fn.smallerThan("sm")]: {
      display: "none",
    },
  },

  social: {
    width: rem(260),

    [theme.fn.smallerThan("sm")]: {
      width: "auto",
      marginLeft: "auto",
    },
  },

  burger: {
    marginRight: theme.spacing.md,

    [theme.fn.largerThan("sm")]: {
      display: "none",
    },
  },

  link: {
    display: "block",
    lineHeight: 1,
    padding: `${rem(8)} ${rem(12)}`,
    borderRadius: theme.radius.sm,
    textDecoration: "none",
    color:
      theme.colorScheme === "dark"
        ? theme.colors.dark[0]
        : theme.colors.gray[7],
    fontSize: theme.fontSizes.sm,
    fontWeight: 500,

    "&:hover": {
      backgroundColor:
        theme.colorScheme === "dark"
          ? theme.colors.dark[6]
          : theme.colors.gray[0],
    },
  },

  linkActive: {
    "&, &:hover": {
      backgroundColor: theme.fn.variant({
        variant: "light",
        color: theme.primaryColor,
      }).background,
      color: theme.fn.variant({ variant: "light", color: theme.primaryColor })
        .color,
    },
  },
}));

interface HeaderMiddleProps {
  links: { link: string; label: string }[];
}

export function HeaderMiddle({ links }: HeaderMiddleProps) {
  const [opened, { toggle }] = useDisclosure(false);
  const [active, setActive] = useState(links[0].link);
  const { classes, cx } = useStyles();
  const myData = useGetMyData();
  const is = useIsLoggedIn();

  const logout = useLogout();

  const items = links.map((link) => {
    if (link.label === "Dashboard" && !is) return null;
    if (link.label === "Profile" && !is) return null;
    if (link.label === "Login" && is) return null;

    if (link.label === "Register" && is) {
      return <Button onClick={() => logout()}>Logout</Button>;
    }

    if (link.label === "Profile" && is)
      return (
        <Link key={link.label} to={link.link} className={cx(classes.link, {})}>
          <Avatar
            size={30}
            src={myData.data?.avatar}
            radius={30}
            style={{ display: "inline-flex" }}
          />
        </Link>
      );

    return (
      <Link key={link.label} to={link.link} className={cx(classes.link, {})}>
        {link.label}
      </Link>
    );
  });

  return (
    <Header height={56} mb={120}>
      <Container className={classes.inner}>
        <Burger
          opened={opened}
          onClick={toggle}
          size="sm"
          className={classes.burger}
        />
        <Group className={classes.links} spacing={5}>
          {items}
        </Group>
        <h1>
          B<span style={{ color: "#f56565" }}>H</span>
        </h1>
      </Container>
    </Header>
  );
}