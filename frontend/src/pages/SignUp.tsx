import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { Link } from "react-router-dom"
import { z } from "zod"

import { AuthCard } from "@/components/AuthCard"
import {
  PasswordStrength,
  passwordScore,
} from "@/components/PasswordStrength"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

const schema = z.object({
  name: z.string().min(1, "what should we call you?"),
  email: z.string().email("hmm, that doesn't look like an email"),
  password: z.string().min(8, "at least 8 characters"),
  terms: z
    .boolean()
    .refine((v) => v === true, "please accept the terms to continue"),
})

type FormValues = z.infer<typeof schema>

const labelClass =
  "font-mono text-[11px] uppercase tracking-wider text-ink-mute"

export function SignUp() {
  const form = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { name: "", email: "", password: "", terms: false },
  })

  const password = form.watch("password")
  const score = passwordScore(password ?? "")

  function onSubmit(values: FormValues) {
    // submit wiring lands in Phase 10
    console.log("sign up", values)
  }

  return (
    <AuthCard>
      <div>
        <h1 className="text-[2rem] leading-tight">create account</h1>
        <p className="mt-1 text-sm text-ink-mute font-hand">
          takes about 30 seconds.
        </p>
      </div>

      <Form {...form}>
        <form
          onSubmit={form.handleSubmit(onSubmit)}
          className="flex flex-col gap-4"
        >
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem className="gap-1.5">
                <FormLabel className={labelClass}>name</FormLabel>
                <FormControl>
                  <Input
                    placeholder="your name"
                    autoComplete="name"
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="email"
            render={({ field }) => (
              <FormItem className="gap-1.5">
                <FormLabel className={labelClass}>email</FormLabel>
                <FormControl>
                  <Input
                    placeholder="you@somewhere.com"
                    autoComplete="email"
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="password"
            render={({ field }) => (
              <FormItem className="gap-1.5">
                <FormLabel className={labelClass}>password</FormLabel>
                <FormControl>
                  <Input
                    type="password"
                    placeholder="• • • • • • •"
                    autoComplete="new-password"
                    {...field}
                  />
                </FormControl>
                <PasswordStrength score={score} className="mt-1" />
                <FormMessage />
              </FormItem>
            )}
          />

          <div className="flex items-start gap-2 text-xs text-ink-mute">
            <Checkbox
              id="terms"
              className="mt-0.5"
              checked={form.watch("terms")}
              onCheckedChange={(v) =>
                form.setValue("terms", v === true, {
                  shouldDirty: true,
                  shouldValidate: true,
                })
              }
            />
            <Label
              htmlFor="terms"
              className="font-hand cursor-pointer leading-snug"
            >
              i agree to the{" "}
              <span className="link-dashed text-ink">terms</span> &{" "}
              <span className="link-dashed text-ink">privacy</span>
            </Label>
          </div>
          {form.formState.errors.terms?.message && (
            <p className="text-sm text-destructive -mt-2">
              {form.formState.errors.terms.message}
            </p>
          )}

          <Button type="submit" className="mt-1">
            create account
          </Button>
        </form>
      </Form>

      <p className="text-center text-sm text-ink-mute font-hand mt-auto">
        have one?{" "}
        <Link to="/sign-in" className="link-dashed text-ink">
          sign in
        </Link>
      </p>
    </AuthCard>
  )
}
