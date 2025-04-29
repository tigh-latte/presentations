return function(input)
	local f = io.open(input.path, "r")
	if not f then error("failed to open file") end

	local lines = f:read("*a")
	f:close()

	print(lines)

	if input.transform then
		local cmd = ""
		if type(input.transform) == "table" then
			for _, part in ipairs(input.transform) do
				cmd = cmd .. " " .. part
			end
		else
			cmd = input.transform
		end

		local proc = io.popen(cmd, "w")
		if not proc then error "failed to open cmd" end
		local _, err = proc:write(lines)
		print("error", err)
		proc:flush()

		lines = proc:read("*a")
		proc:close()
		print(lines)
	end

	local ret = "```" .. (input.lang or "") .. "\n" .. lines .. "\n```"
	return ret
end
