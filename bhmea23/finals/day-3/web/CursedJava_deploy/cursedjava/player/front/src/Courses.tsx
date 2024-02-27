import {
  createStyles,
  Paper,
  Text,
  Title,
  Button,
  rem,
  SimpleGrid,
} from "@mantine/core";

const useStyles = createStyles((theme) => ({
  card: {
    height: rem(440),
    display: "flex",
    flexDirection: "column",
    justifyContent: "space-between",
    alignItems: "flex-start",
    backgroundSize: "cover",
    backgroundPosition: "center",
  },

  title: {
    fontFamily: `Greycliff CF ${theme.fontFamily}`,
    fontWeight: 900,
    color: theme.white,
    lineHeight: 1.2,
    fontSize: rem(32),
    marginTop: theme.spacing.xs,
  },

  category: {
    color: theme.white,
    opacity: 0.7,
    fontWeight: 700,
    textTransform: "uppercase",
  },
}));

interface ArticleCardImageProps {
  image: string;
  title: string;
  category: string;
}

const courses = [
  {
    image:
      "https://images.unsplash.com/photo-1508193638397-1c4234db14d8?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=400&q=80",
    title:
      "Full-Stack Web Development with Java: Build Modern Web Apps from Scratch!",
    category: "Web Development",
  },
  {
    image:
      "https://images.unsplash.com/photo-1508193638397-1c4234db14d8?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=400&q=80",
    title:
      "Mastering gRPC with Java: Build High-Performance Microservices and APIs",
    category: "API Development",
  },
];

export function Courses() {
  const { classes } = useStyles();

  return (
    <SimpleGrid cols={2}>
      {courses.map((course) => (
        <Paper
          shadow="md"
          p="xl"
          radius="md"
          sx={{ backgroundImage: `url(${course.image})` }}
          className={classes.card}
        >
          <div>
            <Text className={classes.category} size="xs">
              {course.category}
            </Text>
            <Title order={3} className={classes.title}>
              {course.title}
            </Title>
          </div>
          <Button disabled variant="white" color="dark">
            Subscribe to Unlock
          </Button>
        </Paper>
      ))}
    </SimpleGrid>
  );
}
