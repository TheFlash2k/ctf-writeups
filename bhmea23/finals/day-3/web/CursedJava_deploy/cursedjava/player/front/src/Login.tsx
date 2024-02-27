import {
  Paper,
  createStyles,
  TextInput,
  PasswordInput,
  Checkbox,
  Button,
  Title,
  Text,
  Anchor,
  rem,
} from '@mantine/core';
import { useLogin } from "./api";

const useStyles = createStyles((theme) => ({
    wrapper:{},
  form: {
    borderRight: `${rem(1)} solid ${theme.colorScheme === "dark" ? theme.colors.dark[7] : theme.colors.gray[3]
      }`,
    minHeight: rem(900),
    maxWidth: rem(450),
    paddingTop: rem(80),

    [theme.fn.smallerThan("sm")]: {
      maxWidth: "100%",
    },
  },

  title: {
    color: theme.colorScheme === "dark" ? theme.white : theme.black,
    fontFamily: `Greycliff CF, ${theme.fontFamily}`,
  },
}));

export function Login() {
  const { classes } = useStyles();

  const [submit, error] = useLogin();

  return (
    <form
      className={classes.wrapper}
      onSubmit={(e: any) => {
        e.preventDefault();
        submit({
          username: e.target.username.value,
          password: e.target.password.value,
        });
      }}
    >
      <Paper className={classes.form} radius={0} p={30}>
        <Title order={2} className={classes.title} ta="center" mt="md" mb={50}>
          Welcome back !
        </Title>

        <TextInput
          label="Username"
          placeholder="Username"
          size="md"
          name="username"
          required
        />
        <PasswordInput
          label="Password"
          placeholder="Password"
          mt="md"
          size="md"
          name="password"
          required
        />
        <span
          style={{
            color: "red",
            marginBottom: "1rem",
          }}
        >
          {error}
        </span>

        <Button fullWidth mt="xl" size="md" type="submit">
          Login
        </Button>
      </Paper>
    </form>
  );
}