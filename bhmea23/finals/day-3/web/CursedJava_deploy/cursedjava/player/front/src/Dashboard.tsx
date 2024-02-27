import {
  createStyles,
  Container,
  Text,
  Button,
  Group,
  rem,
  TextInput,
} from "@mantine/core";
import { useGetFlag, useGetMyData, useRedeem } from "./api";
import { Navigate } from "react-router-dom";
import { Courses } from "./Courses";
const useStyles = createStyles((theme) => ({
  wrapper: {
    position: "relative",
    boxSizing: "border-box",
    backgroundColor:
      theme.colorScheme === "dark" ? theme.colors.dark[8] : theme.white,
  },

  inner: {
    position: "relative",
    paddingTop: rem(200),
    paddingBottom: rem(120),

    [theme.fn.smallerThan("sm")]: {
      paddingBottom: rem(80),
      paddingTop: rem(80),
    },
  },

  title: {
    fontFamily: `Greycliff CF, ${theme.fontFamily}`,
    fontSize: rem(62),
    fontWeight: 900,
    lineHeight: 1.1,
    margin: 0,
    padding: 0,
    color: theme.colorScheme === "dark" ? theme.white : theme.black,

    [theme.fn.smallerThan("sm")]: {
      fontSize: rem(42),
      lineHeight: 1.2,
    },
  },

  description: {
    marginTop: theme.spacing.xl,
    fontSize: rem(24),

    [theme.fn.smallerThan("sm")]: {
      fontSize: rem(18),
    },
  },

  controls: {
    marginTop: `calc(${theme.spacing.xl} * 2)`,

    [theme.fn.smallerThan("sm")]: {
      marginTop: theme.spacing.xl,
    },
  },

  control: {
    height: rem(54),
    paddingLeft: rem(38),
    paddingRight: rem(38),

    [theme.fn.smallerThan("sm")]: {
      height: rem(54),
      paddingLeft: rem(18),
      paddingRight: rem(18),
      flex: 1,
    },
  },
}));

export function Dashboard() {
  const { classes } = useStyles();
  const flag = useGetFlag();
  const mine = useGetMyData();

  const [redeem, error] = useRedeem();

  if (flag.isLoading || mine.isLoading)
    return <div className={classes.wrapper}>loading..</div>;

  if (!mine.data) return <Navigate to="/register" />;

  return (
    <div className={classes.wrapper}>
      <Container size={700} className={classes.inner}>
        <h1 className={classes.title}>
          Welcome{" "}
          <Text
            component="span"
            variant="gradient"
            gradient={{ from: "red", to: "black" }}
            inherit
          >
            {flag.data}
          </Text>{" "}
        </h1>

        <Text className={classes.description} color="black">
          Unfortunately, the website is not ready yet, So the content is not
          available. We will be releasing free Coupons for our first 100
          customers, once our internal testing is done. For now, you can check
          our landing page for more details.
        </Text>

        <Group className={classes.controls}>
          <form
            onSubmit={(e: any) => {
              e.preventDefault();
              redeem(e.target.code.value);
            }}
          >
            <TextInput
              label="Coupon Code"
              placeholder="COUPON CODE"
              size="md"
              name="code"
              required
              style={{
                marginBottom: "1rem",
              }}
            />

            <span
              style={{
                marginBottom: "1rem",
                display: "block",
              }}
            >
              Coupon features are not available right now, we are working on
              fixing the issue.. Get back later
            </span>

            <span
              style={{
                color: "green",
                marginBottom: "1rem",
                display: "block",
              }}
            >
              {error}
            </span>

            <Button
              type="submit"
              size="xl"
              className={classes.control}
              variant="gradient"
              gradient={{ from: "red", to: "black" }}
            >
              Get started
            </Button>
          </form>
        </Group>
      </Container>
    
    </div>
  );
}
