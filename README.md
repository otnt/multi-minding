# Multi-minding

## Background

Multi-minding is a tool that helps programmers to do multiple tasks at the same time.

Based on my personal experience, it is common to get lost when switching among more than three tasks. Especially when each tasks requires non-trivial mind effort. Multi-minding release mind pressure of memorizing what was the state of a task, what task is highest priority etc. You can think of it as the operating system and the human as CPU.

## Overview

On a high level, multi-minding will have several components:

### Task Store

Task store persists all tasks with both current state and history states. It is the source of truth of state of tasks and is the only stateful component. As a result, task store should save data in a decent storage system such as Amazon S3.

### Scheduler

Scheduler decides what task to do given a user request. As an overview, it should consider things like priority, is task blocked, is user qualified etc.

Also, scheduler is the only component modifies task store. When scheduler receives new task (see Task Creator) it will first insert new task to task store, then dispatch the task to user at a proper time (see Task Manager). When task state changed (see Task Manager), it will update task store.

### Task Creator

Task creator creates tasks and send new tasks to scheduler. The goal of task creator is to support general interface for creating tasks using different ways. For example, a user could send an SMS text. Or a bot crawling some news website. Or a monitor watching buglist used in most tech companies.

### Task Manager

Task manager talks with scheduler to fetch tasks to do and notify tasks state. When fetching tasks, it provides current state of user for scheduler to decide what task to dispatch, e.g. where is the user, is user tired, does user have a laptop at hand etc. Alternatively, user could ask for a specific task to do.

When notifying tasks state, it could be a task is finished, is blocked, is de-prioritized etc. Despite the state itself, when notifying a state, it usually means a user is either temporarily or permanently done with a task, and is either ready for another task or needs rest.

### Reminder

As a human, it is inevictably to forget something or just too lazy to fetch a task. Therefore, multi-minding provide reminder as a first-class citizen to remind user of user scheduled tasks and high priority tasks. Reminder repeatedly ask scheduler what are the tasks need to be reminded and has not been done/work in progress, and send remind to user at proper time.

If user ignores the remind, then reminder will simply retry the remind after some time. If user re-acts the remind, it could be switch to remind task, skip the task, delay the task eta. Whatever is user's action, the reminder is simply a relay for task manager. i.e. it only call task manager to finish the work.

### Reporter

Since the goal of multi-minding is to help user better managing multi-tasking, it is beneficial to create a report for user as feedback. The report can include at what time, at what location, and how long is a user doing what task.

Task store will contain all information that reporter needs. Therefore, reporter only needs to query task store to provide feedback.
